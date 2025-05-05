import React, { useState } from "react";
import axios from "axios";

const S3Uploader = () => {
    const [file, setFile] = useState(null);
    const [uploading, setUploading] = useState(false);
    const [message, setMessage] = useState("");

    const handleFileChange = (e) => {
        setFile(e.target.files[0]);
        setMessage("");
    };

    const handleUpload = async () => {
        if (!file) {
            setMessage("Please select a file first.");
            return;
        }

        try {
            setUploading(true);
            setMessage("Getting upload URL...");

            // 1. Get presigned URL from your Go backend
            const filename = encodeURIComponent(file.name); // Optional: sanitize
            const response = await axios.get(
                `http://localhost:8080/api/upload-url?filename=${filename}`
            );
            const { url } = response.data;

            setMessage("Uploading to S3...");

            // 2. Upload file directly to S3
            await axios.put(url, file, {
                headers: {
                    "Content-Type": file.type,
                },
            });

            setMessage("✅ Upload successful!");
        } catch (error) {
            console.error("Upload error:", error);
            setMessage("❌ Upload failed.");
        } finally {
            setUploading(false);
        }
    };

    return (
        <div className="p-4 border rounded-lg shadow-md max-w-md mx-auto">
            <h2 className="text-xl font-bold mb-4">Upload File to S3</h2>
            <input type="file" onChange={handleFileChange} />
            <button
                onClick={handleUpload}
                disabled={!file || uploading}
                className="mt-3 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
                {uploading ? "Uploading..." : "Upload"}
            </button>
            {message && <p className="mt-2 text-sm">{message}</p>}
        </div>
    );
};

export default S3Uploader;
