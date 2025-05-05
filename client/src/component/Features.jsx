import React from 'react';
import { Calendar, Shield, Clock } from 'lucide-react';

const FlexibleFeatures = () => {
    const features = [
        {
            icon: <Calendar className="h-12 w-12 text-black mx-auto mb-4" aria-label="Flexible dates icon" />,
            title: 'Flexible Dates',
            description: 'Choose your check-in and check-out dates with ease, adjusting to fit your schedule.',
        },
        {
            icon: <Shield className="h-12 w-12 text-black mx-auto mb-4" aria-label="Secure booking icon" />,
            title: 'Secure Booking',
            description: 'Book with confidence knowing your reservation is protected with flexible cancellation policies.',
        },
        {
            icon: <Clock className="h-12 w-12 text-black mx-auto mb-4" aria-label="Easy changes icon" />,
            title: 'Easy Changes',
            description: 'Modify or cancel your booking effortlessly if your plans change, with minimal hassle.',
        },
    ];

    return (
        <section className="bg-white font-sans py-12">

                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
                    {features.map((feature, index) => (
                        <article key={index} className="p-6 text-center">
                            {feature.icon}
                            <h3 className="text-xl font-bold text-black mb-2">{feature.title}</h3>
                            <p className="text-base text-gray-600 font-normal">{feature.description}</p>
                        </article>
                    ))}
                </div>

        </section>
    );
};

export default FlexibleFeatures;