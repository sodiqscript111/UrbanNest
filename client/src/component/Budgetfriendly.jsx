import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import PropertyCard from '../component/Propertycard.jsx';

const BudgetFriendly = () => {
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        async function fetchListings() {
            const url = 'https://urbannest-ybda.onrender.com/listings';
            console.log('Fetching Budget-Friendly Nests from:', url);
            try {
                const res = await fetch(url, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });
                console.log('Response status:', res.status);
                if (!res.ok) {
                    const text = await res.text();
                    console.log('Response body:', text);
                    throw new Error(`HTTP error! status: ${res.status}`);
                }
                const data = await res.json();
                console.log('Response data:', data);
                const filteredListings = (data.listings || [])
                    .sort((a, b) => a.price - b.price)
                    .slice(0, 4);
                setListings(filteredListings);
                setLoading(false);
            } catch (err) {
                console.error('Failed to fetch budget-friendly nests:', err);
                setError(`Failed to load budget-friendly nests: ${err.message}`);
                setLoading(false);
            }
        }

        fetchListings();
    }, []);

    const handleExplore = () => {
        navigate('/listings');
        console.log('Explore budget-friendly nests clicked');
    };

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
                    Budget-Friendly Nests in Lagos
                </h2>
                <p className="text-lg font-medium text-gray-600 mb-8">
                    Discover affordable stays in Lagos with great value, comfort, and guest satisfaction.
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
                        No budget-friendly nests available at the moment. Please check back later.
                    </p>
                )}
                <div className="mt-8 text-center">
                    <button
                        onClick={handleExplore}
                        className="py-2 px-4 bg-black text-white font-medium rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                        aria-label="Explore more budget-friendly nests"
                    >
                        Explore
                    </button>
                </div>
            </div>
        </section>
    );
};

export default BudgetFriendly;