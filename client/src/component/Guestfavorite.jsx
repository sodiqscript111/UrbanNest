import React, { useEffect, useState } from 'react';
import PropertyCard from '../component/Propertycard.jsx';

const GuestFavorites = () => {
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchListings() {
            try {
                const res = await fetch('https://urbannest-backend.onrender.com/api/listings', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });
                if (!res.ok) throw new Error('Failed to fetch');
                const data = await res.json();
                // Filter top 4 listings by price (descending)
                const filteredListings = (data.listings || [])
                    .sort((a, b) => b.price - a.price)
                    .slice(0, 4);
                setListings(filteredListings);
                setLoading(false);
            } catch (err) {
                console.error('Failed to fetch guest favorites:', err);
                setError('Failed to load guest favorites. Please try again later.');
                setLoading(false);
            }
        }

        fetchListings();
    }, []);

    if (loading) {
        return (
            <section className="bg-white font-sans py-12">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <p className="text-lg text-gray-600">Loading...</p>
                </div>
            </section>
        );
    }

    if (error) {
        return (
            <section className="bg-white font-sans py-12">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <p className="text-lg text-red-500">{error}</p>
                </div>
            </section>
        );
    }

    return (
        <section className="bg-white font-sans py-12">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <h2 className="text-3xl font-bold text-black mb-4">
                    Guest Favorite Nests in Lagos
                </h2>
                <p className="text-lg font-medium text-gray-600 mb-8">
                    These Nests are highly rated and cherished for their outstanding reviews, reliability, and guest satisfaction.
                </p>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
                    {listings.map((nest) => (
                        <PropertyCard
                            key={nest.id}
                            id={nest.id}
                            media={
                                nest.media && nest.media.length > 0
                                    ? nest.media.map(
                                        (item) => `https://urbannestbucket.s3.eu-north-1.amazonaws.com/${item.media_url}`
                                    )
                                    : ['https://via.placeholder.com/800x600']
                            }
                            name={nest.title}
                            rating="4.5 ★ (100)"
                            description={nest.description}
                            availability={nest.is_available ? 'Available' : 'Booked'}
                            price={nest.price}
                        />
                    ))}
                </div>
                {listings.length === 0 && (
                    <p className="mt-8 text-lg text-gray-600 text-center">
                        No guest favorites available at the moment. Please check back later.
                    </p>
                )}
            </div>
        </section>
    );
};

export default GuestFavorites;