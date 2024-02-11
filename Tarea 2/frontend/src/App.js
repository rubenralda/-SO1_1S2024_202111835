import React, { useState } from 'react';

const CameraCapture = () => {
  const [image, setImage] = useState(null);

  const captureImage = async () => {
    const videoStream = await navigator.mediaDevices.getUserMedia({ video: true });
    const videoElement = document.createElement('video');
    videoElement.srcObject = videoStream;
    videoElement.play();
    const canvas = document.createElement('canvas');
    canvas.width = videoElement.videoWidth;
    canvas.height = videoElement.videoHeight;
    const context = canvas.getContext('2d');
    context.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
    const imageData = canvas.toDataURL('image/png');
    setImage(imageData);
    videoStream.getVideoTracks()[0].stop();
  };

  const uploadImage = async () => {
    try {
      await fetch('/api/upload', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ image: image }),
      });
      console.log('Image uploaded successfully');
    } catch (error) {
      console.error('Error uploading image:', error);
    }
  };

  return (
    <div>
      <button onClick={captureImage}>Capture Image</button>
      {image && <img src={image} alt="Captured" />}
      {image && <button onClick={uploadImage}>Upload Image</button>}
    </div>
  );
};

export default CameraCapture;
