import { useState } from "react";
import axios from "axios";

export default function CreateListing() {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [price, setPrice] = useState("");
    const [location, setLocation] = useState("");
    const [files, setFiles] = useState([]);
    const [uploading, setUploading] = useState(false);
    const [uploadMessages, setUploadMessages] = useState({});
    const [error, setError] = useState(null);

    const BACKEND_URL = "https://urbannest-ybda.onrender.com";
    const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNpbWlkdXNvZGlxdEBnbWFpbC5jb20iLCJleHAiOjE3NDY4ODIyNTgsInVzZXJfaWQiOjF9.SPUCuvUh_2o9Fa1ki80ptgLkh-xnGP6KqGpE89DE3G4";

    const handleFileChange = (e) => {
        const selectedFiles = Array.from(e.target.files);
        if (selectedFiles.length > 10) {
            setError("Maximum 10 files allowed.");
            return;
        }
        if (selectedFiles.some((file) => !["image/jpeg", "image/png"].includes(file.type))) {
            setError("Only JPEG or PNG files allowed.");
            return;
        }
        setFiles(selectedFiles);
        setUploadMessages({});
        setError(null);
    };

    const listProperty = async (e) => {
        e.preventDefault();

        if (!title || !description || !price || !location) {
            setError("Please fill in all listing details.");
            return;
        }

        if (files.length === 0) {
            setError("Please select at least one image.");
            return;
        }

        try {
            setUploading(true);
            setError(null);
            setUploadMessages({});

            const media = [];
            for (const file of files) {
                const filename = encodeURIComponent(file.name);
                setUploadMessages((prev) => ({ ...prev, [filename]: "Getting upload URL..." }));
                console.log(`Requesting presigned URL for ${filename}`);

                const response = await axios.get(
                    `${BACKEND_URL}/api/upload-url?filename=${filename}`,
                    {
                        headers: {
                            Authorization: `Bearer ${TOKEN}`,
                        },
                    }
                );
                const { url } = response.data;
                console.log(`Presigned URL: ${url}`);

                setUploadMessages((prev) => ({ ...prev, [filename]: "Uploading to S3..." }));

                const uploadResponse = await axios.put(url, file, {
                    headers: {
                        "Content-Type": file.type,
                    },
                });

                if (uploadResponse.status !== 200) {
                    throw new Error(`S3 upload failed: Status ${uploadResponse.status}`);
                }
                console.log(`Successfully uploaded ${filename} to S3`);
                setUploadMessages((prev) => ({ ...prev, [filename]: "Upload successful!" }));

                media.push({
                    media_url: filename,
                    media_type: file.type,
                });
            }

            const payload = {
                title,
                description,
                price: parseFloat(price),
                location,
                is_available: true,
                media,
            };
            console.log("Sending payload to /listings:", JSON.stringify(payload, null, 2));

            const res = await axios.post(`${BACKEND_URL}/listings`, payload, {
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${TOKEN}`,
                },
            });

            console.log("Listing created response:", res.data);
            alert("Listing created successfully!");
            setTitle("");
            setDescription("");
            setPrice("");
            setLocation("");
            setFiles([]);
            setUploadMessages({});

        } catch (err) {
            console.error("Failed to create listing:", err);
            setError(err.response?.data?.error || err.message || "Failed to create listing");
        } finally {
            setUploading(false);
        }
    };

    if (error) return <div className="text-red-500">{error}</div>;

    return (
        <div className="p-4 max-w-7xl mx-auto font-sans">
            <h1 className="text-3xl font-bold text-black mb-4">Create New Listing</h1>
            <form onSubmit={listProperty} className="space-y-4">
                <div>
                    <label htmlFor="title" className="block text-sm font-medium text-black">
                        Title
                    </label>
                    <input
                        type="text"
                        id="title"
                        placeholder="e.g., Cozy Studio in Lekki"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                    />
                </div>
                <div>
                    <label htmlFor="description" className="block text-sm font-medium text-black">
                        Description
                    </label>
                    <textarea
                        id="description"
                        placeholder="Describe your property..."
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                    />
                </div>
                <div>
                    <label htmlFor="price" className="block text-sm font-medium text-black">
                        Price (₦/night)
                    </label>
                    <input
                        type="number"
                        id="price"
                        placeholder="e.g., 100000"
                        value={price}
                        onChange={(e) => setPrice(e.target.value)}
                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                    />
                </div>
                <div>
                    <label htmlFor="location" className="block text-sm font-medium text-black">
                        Location
                    </label>
                    <input
                        type="text"
                        id="location"
                        placeholder="e.g., Lekki, Lagos"
                        value={location}
                        onChange={(e) => setLocation(e.target.value)}
                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                    />
                </div>
                <div>
                    <label htmlFor="files" className="block text-sm font-medium text-black">
                        Images (JPEG/PNG, max 10)
                    </label>
                    <input
                        type="file"
                        id="files"
                        multiple
                        accept="image/jpeg,image/png"
                        onChange={handleFileChange}
                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black focus:outline-none focus:ring-2 focus:ring-gray-500"
                    />
                </div>
                {Object.entries(uploadMessages).map(([filename, message]) => (
                    <p key={filename} className="text-sm text-gray-600">{`${filename}: ${message}`}</p>
                ))}
                <button
                    type="submit"
                    disabled={uploading || !title || !description || !price || !location || files.length === 0}
                    className={`w-full py-2 px-4 text-white font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 ${
                        uploading || !title || !description || !price || !location || files.length === 0
                            ? "bg-gray-400 cursor-not-allowed"
                            : "bg-black hover:bg-gray-900"
                    }`}
                >
                    {uploading ? "Uploading..." : "Create Listing"}
                </button>
            </form>
        </div>
    );
}