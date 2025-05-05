import React from 'react';
import PropertyCard from './PropertyCard';

const GuestFavorites = () => {
    const favorites = [
        {
            id: 1,
            media: ['https://images.unsplash.com/photo-1613977257363-707ba9348227?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80'],
            name: 'Lekki Luxury Loft',
            rating: '4.8 ★ (233)',
            description: 'Spacious loft with modern amenities, perfect for a relaxing stay in the heart of Lekki.',
            availability: 'Available',
            price: '250,000',
        },
        {
            id: 2,
            media: ['https://images.unsplash.com/photo-1560448204-e02f11c3d0e2?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80'],
            name: 'Victoria Island Retreat',
            rating: '4.9 ★ (189)',
            description: 'Cozy apartment with stunning ocean views, ideal for business or leisure.',
            availability: 'Booked',
            price: '300,000',
        },
        {
            id: 3,
            media: ['https://i.imgur.com/xfWhBCq.jpeg'],
            name: 'Ikeja Modern Studio',
            rating: '4.7 ★ (156)',
            description: 'Chic studio with fast Wi-Fi and easy access to Lagos’s vibrant city center.',
            availability: 'Available',
            price: '200,000',
        },
        {
            id: 4,
            media: ['https://images.unsplash.com/photo-1600585154340-be6161a56a0c?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80'],
            name: 'Banana Island Villa',
            rating: '5.0 ★ (92)',
            description: 'Luxurious villa with private pool, perfect for an exclusive getaway.',
            availability: 'Booked',
            price: '500,000',
        },
    ];

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
                    {favorites.map((nest) => (
                        <PropertyCard
                            key={nest.id}
                            id={nest.id}
                            media={nest.media}
                            name={nest.name}
                            rating={nest.rating}
                            description={nest.description}
                            availability={nest.availability}
                            price={nest.price}
                        />
                    ))}
                </div>
            </div>
        </section>
    );
};

export default GuestFavorites;