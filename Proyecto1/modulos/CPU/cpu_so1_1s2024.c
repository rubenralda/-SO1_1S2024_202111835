#include <linux/module.h> // THIS_MODULE, MODULE_VERSION, ...
#include <linux/init.h>   // module_{init,exit}
#include <linux/proc_fs.h>
#include <linux/sched/signal.h> // for_each_process()
#include <linux/seq_file.h>
#include <linux/fs.h>
#include <linux/sched.h>
#include <linux/kernel_stat.h>
#include <linux/mm.h> // get_mm_rss()
#include <linux/sched/loadavg.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Ruben Ralda");
MODULE_DESCRIPTION("Informacion CPU uso y procesos");
MODULE_VERSION("1.0");

struct task_struct *task;       // sched.h para tareas/procesos
struct task_struct *task_child; // index de tareas secundarias
struct list_head *list;         // lista de cada tareas

static int escribir_a_proc(struct seq_file *file_proc, void *v)
{
    int running = 0;
    int sleeping = 0;
    int zombie = 0;
    int stopped = 0;
    unsigned long rss;
    unsigned long total_ram_pages;
    
    total_ram_pages = totalram_pages();
    if (!total_ram_pages) {
        pr_err("No memory available\n");
        return -EINVAL;
    }
    
    #ifndef CONFIG_MMU
        pr_err("No MMU, cannot calculate RSS.\n");
        return -EINVAL;
    #endif

    seq_printf(file_proc, "{\"processes\":[\n");
    int b = 0;

    for_each_process(task)
    {
        if (task->mm)
        {
            rss = get_mm_rss(task->mm) << PAGE_SHIFT;
        }
        else
        {
            rss = 0;
        }
        if (b == 0)
        {
            seq_printf(file_proc, "{");
            b = 1;
        }
        else
        {
            seq_printf(file_proc, ",{");
        }
        seq_printf(file_proc, "\"pid\":%d,\n", task->pid);
        seq_printf(file_proc, "\"name\":\"%s\",\n", task->comm);
        seq_printf(file_proc, "\"user\": %d,\n", task->cred->uid.val);
        seq_printf(file_proc, "\"state\":%d,\n", task->__state);
        
        struct task_struct *task_padre;
        task_padre = task->real_parent;
        while (task_padre != NULL)
        {
            if (task_padre->__state == 8193) //en dos procesos
            {
                task_padre = task_padre->real_parent;
            }
            else 
            {
                seq_printf(file_proc, "\"parent\":%d,\n", task_padre->pid);
                task_padre = NULL;
            }
        }   
        seq_printf(file_proc, "\"t\":%d,\n", task->tgid);
        int porcentaje = (rss * 100) / total_ram_pages;
        seq_printf(file_proc, "\"ram\":%d,\n", porcentaje);

        seq_printf(file_proc, "\"child\":[\n");
        int a = 0;
        list_for_each(list, &(task->children))
        {
            task_child = list_entry(list, struct task_struct, sibling);
            if (a != 0)
            {
                seq_printf(file_proc, ",{");
                seq_printf(file_proc, "\"pid\":%d,\n", task_child->pid);
                seq_printf(file_proc, "\"name\":\"%s\",\n", task_child->comm);
                seq_printf(file_proc, "\"state\":%d,\n", task_child->__state);
                seq_printf(file_proc, "\"pidPadre\":%d\n", task->pid);
                seq_printf(file_proc, "}\n");
            }
            else
            {
                seq_printf(file_proc, "{");
                seq_printf(file_proc, "\"pid\":%d,\n", task_child->pid);
                seq_printf(file_proc, "\"name\":\"%s\",\n", task_child->comm);
                seq_printf(file_proc, "\"state\":%d,\n", task_child->__state);
                seq_printf(file_proc, "\"pidPadre\":%d\n", task->pid);
                seq_printf(file_proc, "}\n");
                a = 1;
            }
        }
        a = 0;
        seq_printf(file_proc, "\n]");

        if (task->__state == 0)
        {
            running += 1;
        }
        else if (task->__state == 1)
        {
            sleeping += 1;
        }
        else if (task->__state == 4)
        {
            zombie += 1;
        }
        else
        {
            stopped += 1;
        }
        seq_printf(file_proc, "}\n");
    }
    b = 0;
    seq_printf(file_proc, "],\n");
    /* struct pid *pid_struct;
    pid_t pid_to_find = 996;  // PID del proceso que deseas encontrar

    pid_struct = find_get_pid(pid_to_find);
    if (pid_struct != NULL) {
        task = pid_task(pid_struct, PIDTYPE_PID);
        if (task != NULL) {
            // Proceso encontrado, puedes acceder a su información a través de task
            seq_printf(file_proc, "\"pid\":%d,\n", task->pid);
            seq_printf(file_proc, "\"name\":\"%s\",\n", task->comm);
            seq_printf(file_proc, "\"state\":%ld,\n", task->__state);
            seq_printf(file_proc, "\"pidPadre\":%d,\n", task->parent->pid);
        } else {
            // Proceso no encontrado
        }
    } else {
        // Proceso no encontrado
    } */
    seq_printf(file_proc, "\"running\":%d,\n", running);
    seq_printf(file_proc, "\"sleeping\":%d,\n", sleeping);
    seq_printf(file_proc, "\"zombie\":%d,\n", zombie);
    seq_printf(file_proc, "\"stopped\":%d,\n", stopped);
    seq_printf(file_proc, "\"total\":%d\n", running + sleeping + zombie + stopped);
    seq_printf(file_proc, "}\n");
    return 0;
}

static int escribir_a_cpu(struct seq_file *file_proc, void *v)
{
    seq_printf(file_proc, "%lu.%lu", LOAD_INT(avenrun[0]), LOAD_FRAC(avenrun[0]));
    return 0;
}

static int abrir_aproc_cpu(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_a_cpu, NULL);
}

static int abrir_aproc_procesos(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_a_proc, NULL);
}

static struct proc_ops archivo_operaciones_procesos = {
    .proc_open = abrir_aproc_procesos,
    .proc_read = seq_read
};

static struct proc_ops archivo_operaciones_cpu = {
    .proc_open = abrir_aproc_cpu,
    .proc_read = seq_read
};

static int __init modulo_init(void)
{
    proc_create("cpu_uso", 0, NULL, &archivo_operaciones_cpu);
    proc_create("procesos_cpu", 0, NULL, &archivo_operaciones_procesos);
    printk(KERN_INFO "Insertar Modulo CPU\n");
    return 0;
}

static void __exit modulo_cleanup(void)
{
    remove_proc_entry("procesos_cpu", NULL);
    remove_proc_entry("cpu_uso", NULL);
    printk(KERN_INFO "Remover Modulo CPU\n");
}

module_init(modulo_init);
module_exit(modulo_cleanup);