import redis
import threading

def on_message(message):
    print(f"Mensaje recibido: {message['data'].decode('utf-8')}")

def subscribe_to_channel(channel_name):
    redis_client = redis.Redis(host='10.167.120.43', port=6379)
    pubsub = redis_client.pubsub()
    pubsub.subscribe(**{channel_name: on_message})
    
    # Iniciar un hilo para escuchar mensajes del canal
    thread = pubsub.run_in_thread(sleep_time=0.001)
    print(f"Escuchando el canal '{channel_name}'...")
    
    # Mantener el hilo principal vivo
    try:
        while True:
            pass
    except KeyboardInterrupt:
        print("Saliendo...")
        pubsub.unsubscribe(channel_name)
        thread.stop()

if __name__ == "__main__":
    channel = 'test'
    subscribe_to_channel(channel)
