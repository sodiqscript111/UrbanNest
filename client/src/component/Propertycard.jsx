import React, { useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ChevronLeft, ChevronRight } from 'lucide-react';

const PropertyCard = ({ id, media, name, rating, description, availability, price }) => {
    const scrollRef = useRef(null);
    const [canScrollLeft, setCanScrollLeft] = useState(false);
    const [canScrollRight, setCanScrollRight] = useState(true);
    const navigate = useNavigate();

    if (!id) {
        console.warn('PropertyCard: Missing id prop for property:', name);
    }

    const handleScroll = (direction, e) => {
        e.stopPropagation(); // Prevent card click when clicking arrows
        if (scrollRef.current) {
            const { scrollLeft, clientWidth, scrollWidth } = scrollRef.current;
            const scrollAmount = clientWidth;
            let newScrollLeft;

            if (direction === 'left') {
                newScrollLeft = scrollLeft - scrollAmount;
            } else {
                newScrollLeft = scrollLeft + scrollAmount;
            }

            scrollRef.current.scrollTo({ left: newScrollLeft, behavior: 'smooth' });

            setCanScrollLeft(newScrollLeft > 0);
            setCanScrollRight(newScrollLeft < scrollWidth - clientWidth - 1);
        }
    };

    const handleCardClick = () => {
        if (id) {
            navigate(`/listing/${id}`);
        } else {
            console.error('Cannot navigate: Property id is undefined');
        }
    };

    const images = media && media.length > 0 ? media : ['https://via.placeholder.com/800x600'];

    return (
        <article
            className={`bg-white shadow-md rounded-lg overflow-hidden ${id ? 'cursor-pointer' : 'cursor-default'}`}
            onClick={id ? handleCardClick : undefined}
            aria-label={id ? `View details for ${name}` : `Property ${name} (details unavailable)`}
        >
            <div className="relative group">
                <div
                    ref={scrollRef}
                    className="flex overflow-x-auto scroll-smooth snap-x snap-mandatory hide-scrollbar"
                    style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
                >
                    {images.map((image, index) => (
                        <img
                            key={index}
                            src={image}
                            alt={`${name} - Image ${index + 1}`}
                            className="w-full h-40 sm:h-48 object-cover flex-shrink-0 snap-center"
                            loading="lazy"
                        />
                    ))}
                </div>
                {images.length > 1 && (
                    <>
                        <button
                            onClick={(e) => handleScroll('left', e)}
                            disabled={!canScrollLeft}
                            className={`absolute left-1 sm:left-2 top-1/2 transform -translate-y-1/2 bg-black text-white p-1.5 sm:p-2 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity ${
                                !canScrollLeft ? 'opacity-50 cursor-not-allowed' : ''
                            }`}
                            aria-label="Previous image"
                        >
                            <ChevronLeft className="h-5 w-5 sm:h-6 sm:w-6" />
                        </button>
                        <button
                            onClick={(e) => handleScroll('right', e)}
                            disabled={!canScrollRight}
                            className={`absolute right-1 sm:right-2 top-1/2 transform -translate-y-1/2 bg-black text-white p-1.5 sm:p-2 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity ${
                                !canScrollRight ? 'opacity-50 cursor-not-allowed' : ''
                            }`}
                            aria-label="Next image"
                        >
                            <ChevronRight className="h-5 w-5 sm:h-6 sm:w-6" />
                        </button>
                    </>
                )}
            </div>
            <div className="p-3 sm:p-4">
                <div className="flex justify-between items-center mb-2">
                    <h3 className="text-base sm:text-lg font-bold text-black">{name}</h3>
                    <span className="text-xs sm:text-sm font-medium text-black">{rating}</span>
                </div>
                <p className="text-xs sm:text-sm text-gray-500 mb-2">{description}</p>
                <p
                    className={`text-xs sm:text-sm font-normal mb-2 ${
                        availability === 'Available' ? 'text-green-600' : 'text-red-600'
                    }`}
                >
                    {availability}
                </p>
                <p className="text-sm sm:text-base font-bold text-black">
                    ₦{price.toLocaleString()}/night
                </p>
            </div>
        </article>
    );
};

export default PropertyCard;