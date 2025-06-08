import { UploadPasswordModal } from "./components/UploadModal";
import { encryptFile } from "./utils/cryptoUtils" 
import React, { useState, useRef } from "react";
import { CopyClick } from "./components/CopyClick"
import Header from "./components/Header"

export default function App() {
  const [file, setFile] = useState<File | null>(null);
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [downloadLink, setDownloadLink] = useState("");
  const [encryptedFileName, setEncryptedFileName] = useState("");


  const onFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files?.length) {
      setFile(e.target.files[0]);
    }
  };

  const onUploadClick = () => {
    console.log('Upload button clicked');
    if (!file) return;
    setShowPasswordModal(true);
  };

  const uploadEncryptedFile = async (password: string) => {
    setShowPasswordModal(false);
    if (!file) return;
    setUploading(true);
    try {
      const encryptedBlob = await encryptFile(file, password);

      const formData = new FormData();
      formData.append("file", encryptedBlob, file.name + ".enc");

      const res = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
      });

      const json = await res.json();
      
      function getFileNameFromUrl(url: string): string {
        const pathname = new URL(url).pathname;
        return pathname.substring(pathname.lastIndexOf('/') + 1);
      }
      if (json.download_url) {
        const encFileName = getFileNameFromUrl(json.download_url)
        setEncryptedFileName(encFileName)
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
    <div className="page-container">
      <Header />
      <>
        <div className="container">
          <div className="file-preview-wrapper">
            <label htmlFor="file-upload" className="custom-file-upload">
              Browse
            </label>
            <input 
              id="file-upload"
              type="file" 
              ref={fileInputRef}
              className={isDragOver ? "dragover hidden-input" : "hidden-input"}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onChange={onFileChange}
            />

            {/* Show selected file name, if any */}
            {file && (
              <span className="selected-file-name">{file.name}</span>
            )}

            {file && (
              <button onClick={onUploadClick} disabled={uploading}>
                {uploading ? "Uploading..." : "Encrypt"}
              </button>
            )}


            {showPasswordModal && (
              <UploadPasswordModal
                onConfirm={uploadEncryptedFile}
                onCancel={() => setShowPasswordModal(false)}
              />
            )}

            {downloadLink && (
              <>
                {file && encryptedFileName && (
                  <div className="filename-label-container">
                    <label className="filename-label-og">
                      <strong>{file.name}</strong>
                    </label>
                    <label className="filename-label-arrow"> → </label>
                    <label className="filename-label-enc">
                      <strong>{encryptedFileName}</strong>
                    </label>
                  </div>
                )}

                <div className="mt-4">
                  <CopyClick  downloadUrl={downloadLink}/>
                </div>
              </>
            )}
          </div>
        </div>
      </>
    </div>
  );
}
