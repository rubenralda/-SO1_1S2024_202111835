export namespace main {
	
	export class ramstruct {
	    TotalRam: number;
	    MemoriaEnUso: number;
	    Porcentaje: number;
	    Libre: number;
	
	    static createFrom(source: any = {}) {
	        return new ramstruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TotalRam = source["TotalRam"];
	        this.MemoriaEnUso = source["MemoriaEnUso"];
	        this.Porcentaje = source["Porcentaje"];
	        this.Libre = source["Libre"];
	    }
	}

}

