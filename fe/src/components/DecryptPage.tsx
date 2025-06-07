import { useEffect, useState } from "react";
import { PasswordPromptModal } from "./PasswordPromptModal";
import { deriveKey } from "../utils/cryptoUtils";

export default function DecryptPage() {
  const [fileUrl, setFileUrl] = useState<string | null>(null);
  const [triesLeft, setTriesLeft] = useState(3);
  const [disabled, setDisabled] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const url = params.get("url");
    if (!url) {
      setMessage("Missing file URL");
      return;
    }
    setFileUrl(url);
  }, []);

  const attemptDecrypt = async (password: string) => {
    if (!fileUrl) return;

    try {
      setMessage("Downloading encrypted file...");
      const res = await fetch(fileUrl);
      const encryptedBlob = await res.blob();
      const encryptedArrayBuffer = await encryptedBlob.arrayBuffer();
      const encryptedBytes = new Uint8Array(encryptedArrayBuffer);

      // Parse salt and iv
      const salt = encryptedBytes.slice(0, 16);
      const iv = encryptedBytes.slice(16, 28);
      const data = encryptedBytes.slice(28);

      setMessage("Deriving key...");
      const key = await deriveKey(password, salt);

      setMessage("Decrypting...");
      const decryptedBuffer = await crypto.subtle.decrypt(
        { name: "AES-GCM", iv },
        key,
        data
      );

      // Trigger download
      const decryptedBlob = new Blob([decryptedBuffer]);
      const a = document.createElement("a");
      a.href = URL.createObjectURL(decryptedBlob);

      // Use original filename without ".enc"
      try {
        const originalFilename = fileUrl.split("/").pop() || "download";
        a.download = originalFilename.replace(/\.enc$/, "");
      } catch {
        a.download = "download";
      }

      a.click();
      setMessage("Decryption successful. Download started.");
    } catch (e) {
      const newTries = triesLeft - 1;
      setTriesLeft(newTries);

      if (newTries <= 0) {
        setDisabled(true);
        setMessage("Password incorrect 3 times. Closing page in 5 seconds...");
        setTimeout(() => window.close(), 5000);
      } else {
        setMessage("Password incorrect. Try again.");
      }
    }
  };

  return (
    <div>
      {message && <p>{message}</p>}

      {fileUrl && (
        <PasswordPromptModal
          triesLeft={triesLeft}
          disabled={disabled}
          onSubmit={attemptDecrypt}
        />
      )}
    </div>
  );
}
