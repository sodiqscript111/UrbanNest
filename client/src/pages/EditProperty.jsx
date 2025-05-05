import { useState, useEffect } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";

const EditListing = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [price, setPrice] = useState("");
    const [location, setLocation] = useState("");
    const [isAvailable, setIsAvailable] = useState(true);
    const [files, setFiles] = useState([]);
    const [existingMedia, setExistingMedia] = useState([]);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState(null);
    const [uploadMessages, setUploadMessages] = useState({});

    useEffect(() => {
        async function fetchListing() {
            try {
                const res = await axios.get(`/api/listings/${id}`);
                const listing = res.data.listing;
                setTitle(listing.title);
                setDescription(listing.description);
                setPrice(listing.price.toString());
                setLocation(listing.location);
                setIsAvailable(listing.is_available);
                setExistingMedia(listing.media || []);
            } catch (err) {
                console.error("Failed to fetch listing:", err);
                setError("Failed to fetch listing");
            }
        }
        fetchListing();
    }, [id]);

    const handleFileChange = (e) => {
        const selectedFiles = Array.from(e.target.files);
        if (selectedFiles.length + existingMedia.length > 10) {
            setError("Maximum 10 files allowed.");
            return;
        }
        if (selectedFiles.some(file => !["image/jpeg", "image/png"].includes(file.type))) {
            setError("Only JPEG or PNG files allowed.");
            return;
        }
        setFiles(selectedFiles);
        setUploadMessages({});
        setError(null);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!title || !description || !price || !location) {
            setError("Please fill in all listing details.");
            return;
        }

        try {
            setUploading(true);
            setError(null);
            setUploadMessages({});

            const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NDYzNTMwMzUsInVzZXJfaWQiOjV9.0_yC2AsYhLae0DWqNcE1Zx7HLm0EArGiuaQogjimcEE"; // Replace with dynamic token

            // Upload new files to S3
            const media = [...existingMedia];
            for (const file of files) {
                const filename = encodeURIComponent(file.name);
                setUploadMessages(prev => ({ ...prev, [filename]: "Getting upload URL..." }));
                console.log(`Requesting presigned URL for ${filename}`);

                const response = await axios.get(
                    `http://localhost:8080/api/upload-url?filename=${filename}`
                );
                const { url } = response.data;
                console.log(`Presigned URL: ${url}`);

                setUploadMessages(prev => ({ ...prev, [filename]: "Uploading to S3..." }));

                try {
                    const uploadResponse = await axios.put(url, file, {
                        headers: {
                            "Content-Type": file.type || "application/octet-stream",
                        },
                    });

                    if (uploadResponse.status !== 200) {
                        throw new Error(`S3 upload failed: Status ${uploadResponse.status}`);
                    }
                    console.log(`Successfully uploaded ${filename} to S3`);
                    setUploadMessages(prev => ({ ...prev, [filename]: "✅ Upload successful!" }));

                    media.push({
                        media_url: filename,
                        media_type: file.type.includes("image") ? "image" : "other",
                    });
                } catch (uploadErr) {
                    console.error(`S3 upload error for ${filename}:`, uploadErr);
                    setUploadMessages(prev => ({ ...prev, [filename]: `❌ Upload failed: ${uploadErr.message}` }));
                    throw new Error(`Failed to upload ${filename} to S3: ${uploadErr.message}`);
                }
            }

            // Update listing
            const payload = {
                title,
                description,
                price: parseFloat(price),
                location,
                is_available: isAvailable,
                media,
            };

            await axios.put(`/api/listings/${id}`, payload, {
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
            });

            navigate("/");
        } catch (err) {
            console.error("❌ Failed to update listing:", err);
            setError(err.message || "Failed to update listing");
        } finally {
            setUploading(false);
        }
    };

    const removeExistingMedia = (mediaId) => {
        setExistingMedia(existingMedia.filter(media => media.id !== mediaId));
    };

    if (error) return <div className="text-red-500">{error}</div>;

    return (
        <div className="p-4 max-w-4xl mx-auto">
            <h1 className="text-2xl font-bold mb-4">Edit Listing</h1>
            <form onSubmit={handleSubmit} className="space-y-4">
                <input
                    type="text"
                    placeholder="Title"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    className="w-full p-2 border rounded"
                />
                <textarea
                    placeholder="Description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    className="w-full p-2 border rounded"
                />
                <input
                    type="number"
                    placeholder="Price"
                    value={price}
                    onChange={(e) => setPrice(e.target.value)}
                    className="w-full p-2 border rounded"
                />
                <input
                    type="text"
                    placeholder="Location"
                    value={location}
                    onChange={(e) => setLocation(e.target.value)}
                    className="w-full p-2 border rounded"
                />
                <label className="flex items-center">
                    <input
                        type="checkbox"
                        checked={isAvailable}
                        onChange={(e) => setIsAvailable(e.target.checked)}
                        className="mr-2"
                    />
                    Available
                </label>
                <div>
                    <h3 className="text-lg font-semibold">Current Media</h3>
                    <div className="grid grid-cols-3 gap-2">
                        {existingMedia.map(media => (
                            <div key={media.id} className="relative">
                                <img
                                    src={`https://your-bucket-name.s3.eu-north-1.amazonaws.com/${media.media_url}`}
                                    alt="Media"
                                    className="w-full h-24 object-cover rounded"
                                />
                                <button
                                    type="button"
                                    onClick={() => removeExistingMedia(media.id)}
                                    className="absolute top-0 right-0 bg-red-600 text-white rounded-full w-6 h-6 flex items-center justify-center"
                                >
                                    X
                                </button>
                            </div>
                        ))}
                    </div>
                </div>
                <input
                    type="file"
                    multiple
                    accept="image/jpeg,image/png"
                    onChange={handleFileChange}
                    className="w-full p-2"
                />
                {Object.entries(uploadMessages).map(([filename, message]) => (
                    <p key={filename} className="text-sm text-gray-600">{`${filename}: ${message}`}</p>
                ))}
                <button
                    type="submit"
                    disabled={uploading || !title || !description || !price || !location}
                    className={`w-full p-2 text-white rounded ${
                        uploading || !title || !description || !price || !location
                            ? "bg-gray-400"
                            : "bg-blue-600 hover:bg-blue-700"
                    }`}
                >
                    {uploading ? "Updating..." : "Update Listing"}
                </button>
            </form>
        </div>
    );
};

export default EditListing;