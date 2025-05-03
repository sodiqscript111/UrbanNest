import { useEffect, useState } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";

export default function EditProperty() {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [price, setPrice] = useState("");
    const [location, setLocation] = useState("");
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isUpdating, setIsUpdating] = useState(false);
    const navigate = useNavigate();

    const { id } = useParams();
    useEffect(() => {
        async function fetchListing() {
            setLoading(true);
            try {
                const res = await axios.get(`/api/listings/${id}`);
                setTitle(res.data.title);
                setDescription(res.data.description);
                setPrice(res.data.price);
                setLocation(res.data.location);
                setLoading(false);
            } catch (error) {
                console.log("❌ Failed to fetch listing:", error);
                setError("Failed to load listing.");
                setLoading(false);
            }
        }

        fetchListing();
    }, [id]);

    const handleEditProperty = async (e) => {
        e.preventDefault();
        setIsUpdating(true);

        try {
            const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NDYzNTMwMzUsInVzZXJfaWQiOjV9.0_yC2AsYhLae0DWqNcE1Zx7HLm0EArGiuaQogjimcEE";

            const payload = {
                title,
                description,
                price: parseFloat(price),
                location,
            };
            const res = await axios.put(`/api/listings/${id}`, payload, {
                headers: {
                    "Content-Type": "application/json",
                    "Accept": "application/json",
                    Authorization: `Bearer ${token}`, // Include the token
                },
            });

            console.log("Listing updated:", res.data);
            setIsUpdating(false);
            navigate('/');
        } catch (error) {
            console.log("❌ Failed to update listing:", error);
            setError("Failed to update listing.");
            setIsUpdating(false);
        }
    };

    if (loading) {
        return <div>Loading listing details...</div>;
    }

    if (error) {
        return <div>{error}</div>;
    }

    return (
        <div>
            <h2>Currently Editing:</h2>
            <div>
                <strong>Title:</strong> {title}
            </div>
            <div>
                <strong>Description:</strong> {description}
            </div>
            <div>
                <strong>Price:</strong> {price}
            </div>
            <div>
                <strong>Location:</strong> {location}
            </div>
            <br />
            <form onSubmit={handleEditProperty}> {/* Corrected onSubmit handler name */}
                <input
                    type="text"
                    value={title ?? ""}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="Title"
                />
                <textarea
                    value={description ?? ""}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="Description"
                />
                <input
                    type="number"
                    value={price ?? ""}
                    onChange={(e) => setPrice(e.target.value)}
                    placeholder="Price"
                />
                <input
                    type="text"
                    value={location ?? ""}
                    onChange={(e) => setLocation(e.target.value)}
                    placeholder="Location"
                />
                <button type="submit" disabled={isUpdating}>
                    {isUpdating ? "Updating..." : "Update Property"}
                </button>
            </form>
        </div>
    );
}