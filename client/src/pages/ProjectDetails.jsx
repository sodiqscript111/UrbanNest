import React, { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { ChevronLeft, ChevronRight, X } from 'lucide-react';
import { PaystackButton } from 'react-paystack';

const PropertyDetails = () => {
    const { id } = useParams();
    const [listing, setListing] = useState(null);
    const [error, setError] = useState(null);
    const [checkIn, setCheckIn] = useState(null);
    const [checkOut, setCheckOut] = useState(null);
    const [bookingError, setBookingError] = useState('');
    const [selectedImage, setSelectedImage] = useState(null);
    const [formData, setFormData] = useState({
        email: '',
        name: '',
        phone: '',
    });
    const [paymentStatus, setPaymentStatus] = useState(null);
    const scrollRef = useRef(null);
    const modalRef = useRef(null);
    const [canScrollLeft, setCanScrollLeft] = useState(false);
    const [canScrollRight, setCanScrollRight] = useState(true);

    const publicKey = import.meta.env.VITE_PAYSTACK_PUBLIC_KEY || 'pk_test_your_test_public_key';
    const customerId = 1; // Hardcoded for testing; replace with auth system

    useEffect(() => {
        async function fetchListing() {
            try {
                const res = await axios.get(`/api/listings/${id}`);
                if (!res.data.listing) {
                    throw new Error('Listing not found');
                }
                setListing(res.data.listing);
            } catch (err) {
                console.error('Failed to fetch listing:', err.message, err.response?.data);
                setError('Failed to fetch listing details');
            }
        }
        fetchListing();
    }, [id]);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (modalRef.current && !modalRef.current.contains(event.target)) {
                setSelectedImage(null);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleScroll = (direction) => {
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

    const openImageModal = (image) => {
        setSelectedImage(image);
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const calculateTotalAmount = () => {
        if (!checkIn || !checkOut || !listing || !listing.price) {
            console.warn('Invalid amount calculation inputs:', { checkIn, checkOut, listing });
            return 0;
        }
        const nights = Math.ceil((checkOut - checkIn) / (1000 * 60 * 60 * 24));
        if (nights <= 0 || isNaN(listing.price)) {
            console.warn('Invalid nights or price:', { nights, price: listing.price });
            return 0;
        }
        const amount = Math.round(nights * listing.price * 100); // Convert to kobo
        console.log('Calculated amount (kobo):', amount);
        return amount;
    };

    const createBooking = async (reference) => {
        try {
            const response = await axios.post('/api/booking', {
                listing_id: parseInt(id),
                customer_id: customerId,
                start_date: checkIn.toISOString().split('T')[0],
                end_date: checkOut.toISOString().split('T')[0],
                status: 'confirmed',
                payment_reference: reference,
            });
            setPaymentStatus('success');
            setBookingError('');
            return response.data.booking;
        } catch (err) {
            console.error('Failed to create booking:', err.response?.data || err.message);
            setPaymentStatus('failed');
            setBookingError(err.response?.data?.error || 'Failed to create booking');
            throw err;
        }
    };

    const paystackProps = {
        email: formData.email,
        amount: calculateTotalAmount(),
        publicKey,
        text: 'Book Now',
        metadata: {
            name: String(formData.name),
            phone: String(formData.phone),
            listing_id: String(id),
            customer_id: String(customerId),
        },
        reference: `booking_${id}_${new Date().getTime()}`,
        onSuccess: (reference) => {
            console.log('Paystack success:', reference);
            createBooking(reference.reference)
                .then(() => {
                    alert('Payment and booking successful!');
                    setFormData({ email: '', name: '', phone: '' });
                    setCheckIn(null);
                    setCheckOut(null);
                })
                .catch(() => {
                    alert('Payment succeeded, but booking failed. Please contact support.');
                });
        },
        onClose: () => {
            console.log('Paystack closed');
            setPaymentStatus('cancelled');
            setBookingError('Payment was cancelled');
        },
    };

    const validateEmail = (email) => {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    };

    const handleBooking = (e) => {
        e.preventDefault();
        console.log('Paystack props:', paystackProps);

        if (!formData.email || !validateEmail(formData.email)) {
            setBookingError('Please enter a valid email address');
            return;
        }
        if (!formData.name) {
            setBookingError('Please enter your name');
            return;
        }
        if (!formData.phone || !/^\+?\d{10,15}$/.test(formData.phone)) {
            setBookingError('Please enter a valid phone number (e.g., +2341234567890)');
            return;
        }
        if (!checkIn || !checkOut) {
            setBookingError('Please select check-in and check-out dates');
            return;
        }
        if (checkOut <= checkIn) {
            setBookingError('Check-out must be after check-in');
            return;
        }
        const amount = calculateTotalAmount();
        if (amount <= 0) {
            setBookingError('Invalid booking amount. Please check dates and listing price.');
            return;
        }
        setBookingError('');
        // PaystackButton will handle payment
    };

    if (error) return <div className="text-red-500 p-4 max-w-7xl mx-auto">{error}</div>;
    if (!listing) return <div className="p-4 max-w-7xl mx-auto">Loading...</div>;

    const images = listing.media && listing.media.length > 0
        ? listing.media.map((item) => `https://urbannestbucket.s3.eu-north-1.amazonaws.com/${item.media_url}`)
        : ['https://via.placeholder.com/800x600'];

    return (
        <section className="bg-white font-sans py-12">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="bg-white shadow-lg rounded-2xl p-6 sm:p-8">
                    <h1 className="text-2xl sm:text-3xl font-bold text-black mb-4">{listing.title}</h1>
                    <div className="relative group mb-6">
                        <div
                            ref={scrollRef}
                            className="flex overflow-x-auto scroll-smooth snap-x snap-mandatory hide-scrollbar"
                            style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
                        >
                            {images.map((image, index) => (
                                <img
                                    key={index}
                                    src={image}
                                    alt={`${listing.title} - Image ${index + 1}`}
                                    className="w-full h-80 sm:h-[32rem] object-cover flex-shrink-0 snap-center cursor-pointer"
                                    loading="lazy"
                                    onClick={() => openImageModal(image)}
                                />
                            ))}
                        </div>
                        {images.length > 1 && (
                            <>
                                <button
                                    onClick={() => handleScroll('left')}
                                    disabled={!canScrollLeft}
                                    className={`absolute left-4 sm:left-6 top-1/2 transform -translate-y-1/2 bg-black text-white p-2 sm:p-4 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity ${
                                        !canScrollLeft ? 'opacity-50 cursor-not-allowed' : ''
                                    }`}
                                    aria-label="Previous image"
                                >
                                    <ChevronLeft className="h-8 w-8 sm:h-10 sm:w-10" />
                                </button>
                                <button
                                    onClick={() => handleScroll('right')}
                                    disabled={!canScrollRight}
                                    className={`absolute right-4 sm:right-6 top-1/2 transform -translate-y-1/2 bg-black text-white p-2 sm:p-4 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity ${
                                        !canScrollRight ? 'opacity-50 cursor-not-allowed' : ''
                                    }`}
                                    aria-label="Next image"
                                >
                                    <ChevronRight className="h-8 w-8 sm:h-10 sm:w-10" />
                                </button>
                            </>
                        )}
                    </div>

                    {selectedImage && (
                        <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
                            <div ref={modalRef} className="relative max-w-5xl w-full p-4">
                                <button
                                    onClick={() => setSelectedImage(null)}
                                    className="absolute top-2 right-2 bg-black text-white p-2 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                    aria-label="Close image modal"
                                >
                                    <X className="h-6 w-6" />
                                </button>
                                <img
                                    src={selectedImage}
                                    alt="Full-size property image"
                                    className="w-full max-h-[90vh] object-contain"
                                />
                            </div>
                        </div>
                    )}

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        <div className="md:col-span-2">
                            <div className="flex justify-between items-center mb-4">
                                <h2 className="text-lg sm:text-xl font-bold text-black">{listing.title}</h2>
                                <span className="text-sm sm:text-base font-medium text-black">4.5 ★ (100)</span>
                            </div>
                            <p className="text-sm sm:text-base text-gray-500 mb-4">{listing.description}</p>
                            <p
                                className={`text-sm sm:text-base font-normal mb-4 ${
                                    listing.is_available ? 'text-green-600' : 'text-red-600'
                                }`}
                            >
                                {listing.is_available ? 'Available' : 'Booked'}
                            </p>
                            <p className="text-base sm:text-lg font-bold text-black mb-4">
                                ₦{listing.price.toLocaleString()}/night
                            </p>
                            <p className="text-sm sm:text-base text-gray-600 mb-4">
                                Location: {listing.location}
                            </p>
                        </div>
                        <div className="md:col-span-1">
                            <form onSubmit={handleBooking} className="space-y-4 bg-gray-50 p-4 sm:p-6 rounded-lg">
                                <h3 className="text-lg sm:text-xl font-bold text-black mb-4">Book Your Stay</h3>
                                <div>
                                    <label htmlFor="email" className="block text-sm font-medium text-black">
                                        Email
                                    </label>
                                    <input
                                        id="email"
                                        name="email"
                                        type="email"
                                        value={formData.email}
                                        onChange={handleInputChange}
                                        placeholder="Enter your email"
                                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                        aria-label="Email address"
                                        required
                                    />
                                </div>
                                <div>
                                    <label htmlFor="name" className="block text-sm font-medium text-black">
                                        Name
                                    </label>
                                    <input
                                        id="name"
                                        name="name"
                                        type="text"
                                        value={formData.name}
                                        onChange={handleInputChange}
                                        placeholder="Enter your name"
                                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                        aria-label="Full name"
                                        required
                                    />
                                </div>
                                <div>
                                    <label htmlFor="phone" className="block text-sm font-medium text-black">
                                        Phone
                                    </label>
                                    <input
                                        id="phone"
                                        name="phone"
                                        type="tel"
                                        value={formData.phone}
                                        onChange={handleInputChange}
                                        placeholder="Enter your phone number (e.g., +2341234567890)"
                                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                        aria-label="Phone number"
                                        required
                                    />
                                </div>
                                <div>
                                    <label htmlFor="check-in" className="block text-sm font-medium text-black">
                                        Check-In
                                    </label>
                                    <DatePicker
                                        id="check-in"
                                        selected={checkIn}
                                        onChange={(date) => setCheckIn(date)}
                                        minDate={new Date()}
                                        placeholderText="Select date"
                                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                        aria-label="Check-in date"
                                    />
                                </div>
                                <div>
                                    <label htmlFor="check-out" className="block text-sm font-medium text-black">
                                        Check-Out
                                    </label>
                                    <DatePicker
                                        id="check-out"
                                        selected={checkOut}
                                        onChange={(date) => {
                                            if (checkIn && date <= checkIn) {
                                                setBookingError('Check-out must be after check-in');
                                                return;
                                            }
                                            setCheckOut(date);
                                            setBookingError('');
                                        }}
                                        minDate={checkIn || new Date()}
                                        placeholderText="Select date"
                                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                        aria-label="Check-out date"
                                    />
                                </div>
                                {bookingError && <p className="text-red-500 text-sm">{bookingError}</p>}
                                {paymentStatus === 'success' && (
                                    <p className="text-green-600 text-sm">Booking created successfully!</p>
                                )}
                                {paymentStatus === 'failed' && (
                                    <p className="text-red-500 text-sm">Payment failed. Please try again.</p>
                                )}
                                <PaystackButton
                                    className="w-full py-2 px-4 bg-black text-white font-medium rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                    {...paystackProps}
                                />
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default PropertyDetails;