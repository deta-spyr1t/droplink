import { UploadPasswordModal } from "./components/UploadPasswordModal";
import { encryptFile } from "./utils/cryptoUtils" 
import React, { useState, useRef } from "react";

export default function App() {
  const [file, setFile] = useState<File | null>(null);
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [downloadLink, setDownloadLink] = useState("");

  const onFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files?.length) {
      setFile(e.target.files[0]);
    }
  };

  const onUploadClick = () => {
    if (!file) return;
    setShowPasswordModal(true);
  };

  const uploadEncryptedFile = async (password: string) => {
    setShowPasswordModal(false);
    if (!file) return;
    setUploading(true);
    try {
      const encryptedBlob = await encryptFile(file, password);

      // Create form data and upload
      const formData = new FormData();
      formData.append("file", encryptedBlob, file.name + ".enc");

      const res = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
      });

      const json = await res.json();

      if (json.download_url) {
        const frontendBaseUrl = window.location.origin; // e.g. http://localhost:3000
        const decryptLink = `${frontendBaseUrl}/decrypt?url=${encodeURIComponent(json.download_url)}`;
        setDownloadLink(decryptLink);
      } else {
        alert("Upload failed");
      }
    } catch (e) {
      alert("Encryption/upload failed: " + (e as Error).message);
    }
    setUploading(false);
  };

  const [isDragOver, setIsDragOver] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Handlers
  function handleDragOver(e: React.DragEvent) {
    e.preventDefault();
    setIsDragOver(true);
  }

  function handleDragLeave(e: React.DragEvent) {
    e.preventDefault();
    setIsDragOver(false);
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    setIsDragOver(false);

    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      // Assign dropped files to the file input element
      if (fileInputRef.current) {
        fileInputRef.current.files = e.dataTransfer.files;
      }
      e.dataTransfer.clearData();
      // Optionally do more with the dropped file here...
    }
  }

  return (
    <div className="container">
      <div>
        <p>End-to-End Encrypted file sharing service</p>
        <input 
          type="file" 
          ref={fileInputRef}
          className={isDragOver ? "dragover" : ""}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onChange={onFileChange}
        />
        {file && <p>Selected: {file.name}</p>}
        <button disabled={!file || uploading} onClick={onUploadClick}>
          {uploading ? "Uploading..." : "Encrypt"}
        </button>

        {showPasswordModal && (
          <UploadPasswordModal
            onConfirm={uploadEncryptedFile}
            onCancel={() => setShowPasswordModal(false)}
          />
        )}

        {downloadLink && (
          <p>
            Download Link: <br />
            <a href={downloadLink} target="_blank" rel="noopener noreferrer">
              {downloadLink}
            </a>
          </p>
        )}
      </div>
    </div>
  );
}
