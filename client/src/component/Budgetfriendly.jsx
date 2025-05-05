
import {useNavigate} from "react-router-dom";
import React from 'react';
import PropertyCard from './PropertyCard';

const BudgetFriendly = () => {
    const budgetNests = [
        {
            id: 1,
            media: ['https://images.unsplash.com/photo-1600585154340-be6161a56a0c?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80'],
            name: 'Ajah Cozy Studio',
            rating: '4.6 ★ (120)',
            description: 'Affordable studio with essential amenities, close to Ajah’s bustling markets.',
            availability: 'Available',
            price: '100,000',
        },
        {
            id: 2,
            media: ['https://images.unsplash.com/photo-1560448204-e02f11c3d0e2?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80'],
            name: 'Yaba Budget Apartment',
            rating: '4.5 ★ (98)',
            description: 'Compact apartment perfect for students or travelers, near Yaba’s vibrant scene.',
            availability: 'Booked',
            price: '120,000',
        },
        {
            id: 3,
            media: ['https://i.imgur.com/xfWhBCq.jpeg'],
            name: 'Surulere Modern Flat',
            rating: '4.7 ★ (145)',
            description: 'Stylish flat with great value, ideal for short stays in Surulere.',
            availability: 'Available',
            price: '150,000',
        },
        {
            id: 4,
            media: ['https://images.unsplash.com/photo-1613977257363-707ba9348227?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80'],
            name: 'Ikorodu Simple Loft',
            rating: '4.4 ★ (76)',
            description: 'Minimalist loft with budget-friendly pricing, great for a quiet retreat.',
            availability: 'Booked',
            price: '90,000',
        },
    ];

    const handleExplore = () => {
        console.log('Explore budget-friendly nests clicked');
    };

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
                    {budgetNests.map((nest) => (
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