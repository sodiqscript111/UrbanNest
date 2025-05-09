import PropertyCard from "../component/Propertycard.jsx";
import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Search } from "lucide-react";

const SkeletonLoader = () => {
    return (
        <motion.div
            className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.4 }}
        >
            <div className="h-6 bg-gray-200 rounded w-1/2 mb-6 animate-pulse"></div>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {[...Array(4)].map((_, index) => (
                    <motion.div
                        key={index}
                        className="bg-white rounded-lg shadow-md overflow-hidden animate-pulse"
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.4, delay: index * 0.1 }}
                    >
                        <div className="w-full h-44 bg-gray-200"></div>
                        <div className="p-4">
                            <div className="h-5 bg-gray-200 rounded w-3/4 mb-2"></div>
                            <div className="h-4 bg-gray-200 rounded w-full mb-2"></div>
                            <div className="h-4 bg-gray-200 rounded w-5/6 mb-2"></div>
                            <div className="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
                            <div className="h-5 bg-gray-200 rounded w-1/3"></div>
                        </div>
                    </motion.div>
                ))}
            </div>
        </motion.div>
    );
};

export default function Home() {
    const [listings, setListings] = useState([]);
    const [filteredListings, setFilteredListings] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [searchQuery, setSearchQuery] = useState("");

    // Scroll to top on mount
    useEffect(() => {
        window.scrollTo(0, 0);
    }, []);

    useEffect(() => {
        async function fetchWithRetry(url, options, retries = 3, delay = 1000) {
            for (let i = 0; i < retries; i++) {
                try {
                    console.log(`Attempt ${i + 1} - Fetching listings from:`, url);
                    const res = await fetch(url, options);
                    console.log('Response status:', res.status);
                    if (!res.ok) {
                        const text = await res.text();
                        console.log('Response body:', text);
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    const data = await res.json();
                    console.log('Response data:', data);
                    return data;
                } catch (err) {
                    console.error(`Attempt ${i + 1} failed:`, err);
                    if (i < retries - 1) {
                        console.log(`Retrying after ${delay}ms...`);
                        await new Promise((resolve) => setTimeout(resolve, delay));
                    } else {
                        throw err;
                    }
                }
            }
        }

        async function fetchListings() {
            const url = 'https://urbannest-ybda.onrender.com/listings';
            try {
                const data = await fetchWithRetry(
                    url,
                    {
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                    },
                    3,
                    1000
                );
                const sortedListings = (data.listings || []).sort((a, b) => b.price - a.price);
                setListings(sortedListings);
                setFilteredListings(sortedListings);
            } catch (err) {
                console.error('Failed to fetch listings:', err);
                setError(`Failed to load listings: ${err.message}`);
            } finally {
                setIsLoading(false);
            }
        }

        fetchListings();
    }, []);

    // Handle search
    useEffect(() => {
        if (searchQuery.trim() === "") {
            setFilteredListings(listings);
        } else {
            const query = searchQuery.toLowerCase();
            const filtered = listings.filter(
                (listing) =>
                    listing.title.toLowerCase().includes(query) ||
                    listing.location.toLowerCase().includes(query)
            );
            setFilteredListings(filtered);
        }
    }, [searchQuery, listings]);

    const handleSearchChange = (e) => {
        setSearchQuery(e.target.value);
    };

    const handleSearchSubmit = (e) => {
        e.preventDefault();
        // Trigger search (already handled by useEffect)
    };

    if (isLoading) {
        return <SkeletonLoader />;
    }

    if (error) {
        return (
            <motion.div
                className="text-red-500 p-4 sm:p-6 max-w-7xl mx-auto bg-white rounded-md shadow-md font-inter"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.4 }}
            >
                {error}
            </motion.div>
        );
    }

    return (
        <div className="bg-gray-100 py-10 sm:py-12 font-inter">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="mb-8 sm:mb-10">
                    <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 tracking-tight leading-tight">
                        Discover Your Perfect Stay
                    </h1>
                </div>
                {/* Search Bar */}
                <motion.form
                    onSubmit={handleSearchSubmit}
                    className="mb-6 sm:mb-8 sticky top-0 z-10 bg-gray-100 py-4 -mx-4 px-4"
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.4 }}
                >
                    <div className="flex items-center gap-2">
                        <input
                            type="text"
                            value={searchQuery}
                            onChange={handleSearchChange}
                            placeholder="Search by title or location..."
                            className="flex-1 px-4 py-3 text-base sm:text-lg text-gray-900 placeholder-gray-500 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 transition-all duration-200"
                            aria-label="Search listings"
                        />
                        <motion.button
                            type="submit"
                            className="bg-white border border-black text-black font-semibold px-4 py-3 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 flex items-center gap-2"
                            whileHover={{ scale: 1.05 }}
                            whileTap={{ scale: 0.95 }}
                            transition={{ duration: 0.2 }}
                            aria-label="Search"
                        >
                            <Search className="h-5 w-5" />
                            <span className="hidden sm:inline">Search</span>
                        </motion.button>
                    </div>
                </motion.form>
                <AnimatePresence>
                    <motion.div
                        className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.4 }}
                    >
                        {filteredListings.map((listing, index) => (
                            <motion.div
                                key={listing.id}
                                className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-200"
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ duration: 0.4, delay: index * 0.1 }}
                                whileHover={{ scale: 1.03 }}
                                whileTap={{ scale: 0.95 }}
                            >
                                <PropertyCard
                                    id={listing.id}
                                    media={
                                        listing.media && listing.media.length > 0
                                            ? listing.media.map(
                                                (item) => `https://urbannestbucket.s3.eu-north-1.amazonaws.com/${item.media_url}`
                                            )
                                            : ["https://via.placeholder.com/600x400"]
                                    }
                                    name={listing.title}
                                    rating="4.5 ★ (100)"
                                    description={listing.description.substring(0, 100) + "..."}
                                    availability={listing.is_available ? "Available" : "Booked"}
                                    price={listing.price}
                                    className="rounded-lg"
                                />
                            </motion.div>
                        ))}
                    </motion.div>
                </AnimatePresence>
                {filteredListings.length === 0 && !error && (
                    <motion.div
                        className="mt-8 sm:mt-10 text-center text-gray-600 text-base sm:text-lg font-inter"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ duration: 0.4 }}
                    >
                        <p>No listings match your search. Try a different query.</p>
                    </motion.div>
                )}
            </div>
        </div>
    );
}