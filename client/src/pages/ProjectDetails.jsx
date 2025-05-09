import React, { useEffect, useRef, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { ChevronLeft, ChevronRight, X } from 'lucide-react';
import { PaystackButton } from 'react-paystack';
import { motion } from 'framer-motion';

const SkeletonLoader = () => {
    return (
        <motion.div
            className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
        >
            <div className="bg-white shadow-lg rounded-2xl p-4 sm:p-6">
                <div className="h-8 w-8 bg-gray-200 rounded-full mb-4 animate-pulse"></div>
                <div className="h-6 bg-gray-200 rounded w-3/4 mb-4 animate-pulse"></div>
                <div className="relative mb-6">
                    <div className="flex overflow-x-auto space-x-4">
                        <div className="w-full h-64 sm:h-80 bg-gray-200 rounded-lg flex-shrink-0 animate-pulse"></div>
                        <div className="w-full h-64 sm:h-80 bg-gray-200 rounded-lg flex-shrink-0 animate-pulse"></div>
                    </div>
                </div>
                <div className="space-y-4">
                    <div className="h-5 bg-gray-200 rounded w-1/2 mb-4 animate-pulse"></div>
                    <div className="h-4 bg-gray-200 rounded w-full mb-2 animate-pulse"></div>
                    <div className="h-4 bg-gray-200 rounded w-5/6 mb-2 animate-pulse"></div>
                    <div className="h-4 bg-gray-200 rounded w-1/4 mb-4 animate-pulse"></div>
                    <div className="h-5 bg-gray-200 rounded w-1/3 mb-4 animate-pulse"></div>
                    <div className="h-4 bg-gray-200 rounded w-1/2 mb-4 animate-pulse"></div>
                    <div className="space-y-4 bg-gray-50 p-4 rounded-lg">
                        <div className="h-5 bg-gray-200 rounded w-1/2 mb-4 animate-pulse"></div>
                        <div className="h-10 bg-gray-200 rounded w-full mb-4 animate-pulse"></div>
                        <div className="h-10 bg-gray-200 rounded w-full mb-4 animate-pulse"></div>
                        <div className="h-10 bg-gray-200 rounded w-full mb-4 animate-pulse"></div>
                        <div className="h-10 bg-gray-200 rounded w-full mb-4 animate-pulse"></div>
                        <div className="h-10 bg-gray-200 rounded w-full animate-pulse"></div>
                    </div>
                </div>
            </div>
        </motion.div>
    );
};

const PropertyDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
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

    // Scroll to top on mount or id change
    useEffect(() => {
        window.scrollTo(0, 0);
    }, [id]);

    useEffect(() => {
        async function fetchWithRetry(url, options, retries = 5, delay = 2000) {
            for (let i = 0; i < retries; i++) {
                try {
                    console.log(`Attempt ${i + 1} - Fetching listing from:`, url);
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

        async function fetchListing() {
            try {
                const data = await fetchWithRetry(
                    `https://urbannest-ybda.onrender.com/listings/${id}`,
                    {
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                    },
                    5,
                    2000
                );
                if (!data.listing) {
                    throw new Error('Listing not found');
                }
                setListing(data.listing);
            } catch (err) {
                console.error('Failed to fetch listing:', err.message);
                setError(`Listing not found (ID: ${id})`);
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
            return 0;
        }
        const nights = Math.ceil((checkOut - checkIn) / (1000 * 60 * 60 * 24));
        if (nights <= 0 || isNaN(listing.price)) {
            return 0;
        }
        const amount = Math.round(nights * listing.price * 100); // Convert to kobo
        console.log('Calculated amount (kobo):', amount);
        return amount;
    };

    const createBooking = async (reference) => {
        try {
            const response = await fetch('https://urbannest-ybda.onrender.com/api/booking', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    listing_id: parseInt(id),
                    customer_id: customerId,
                    start_date: checkIn.toISOString().split('T')[0],
                    end_date: checkOut.toISOString().split('T')[0],
                    status: 'confirmed',
                    payment_reference: reference,
                }),
            });
            if (!response.ok) {
                const text = await response.text();
                throw new Error(`HTTP error! status: ${response.status}, body: ${text}`);
            }
            const data = await response.json();
            setPaymentStatus('success');
            setBookingError('');
            return data.booking;
        } catch (err) {
            console.error('Failed to create booking:', err.message);
            setPaymentStatus('failed');
            setBookingError(err.message || 'Failed to create booking');
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

    // Animation variants
    const imageVariants = {
        hidden: { opacity: 0, y: 20 },
        visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: 'easeOut' } },
    };

    const formVariants = {
        hidden: { opacity: 0, y: 20 },
        visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: 'easeOut' } },
    };

    const modalVariants = {
        hidden: { opacity: 0, scale: 0.95 },
        visible: { opacity: 1, scale: 1, transition: { duration: 0.2, ease: 'easeOut' } },
    };

    const detailsVariants = {
        hidden: { opacity: 0, y: 10 },
        visible: { opacity: 1, y: 0, transition: { duration: 0.4, delay: 0.1, ease: 'easeOut' } },
    };

    const messageVariants = {
        hidden: { opacity: 0 },
        visible: { opacity: 1, transition: { duration: 0.2 } },
    };

    if (error) return (
        <motion.div
            className="text-red-500 p-4 max-w-7xl mx-auto"
            variants={messageVariants}
            initial="hidden"
            animate="visible"
        >
            {error}
        </motion.div>
    );
    if (!listing) return <SkeletonLoader />;

    const images = listing.media && listing.media.length > 0
        ? listing.media.map((item) => `https://urbannestbucket.s3.eu-north-1.amazonaws.com/${item.media_url}`)
        : ['https://via.placeholder.com/800x600'];

    return (
        <motion.section
            className="bg-white font-sans py-8"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
        >
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="bg-white shadow-lg rounded-2xl p-4 sm:p-6">
                    <motion.button
                        onClick={() => navigate(-1)}
                        className="mb-4 bg-black text-white p-2.5 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                        aria-label="Go back"
                        whileHover={{ scale: 1.1 }}
                        whileTap={{ scale: 0.9 }}
                        transition={{ duration: 0.2 }}
                    >
                        <ChevronLeft className="h-6 w-6" />
                    </motion.button>
                    <motion.h1
                        className="text-xl sm:text-2xl md:text-3xl font-bold text-black mb-4"
                        variants={detailsVariants}
                        initial="hidden"
                        animate="visible"
                    >
                        {listing.title}
                    </motion.h1>
                    <div className="relative group mb-6">
                        <div
                            ref={scrollRef}
                            className="flex overflow-x-auto scroll-smooth snap-x snap-mandatory hide-scrollbar touch-pan-x"
                            style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
                        >
                            {images.map((image, index) => (
                                <motion.img
                                    key={index}
                                    src={image}
                                    alt={`${listing.title} - Image ${index + 1}`}
                                    className="w-full h-64 sm:h-80 md:h-[32rem] object-cover flex-shrink-0 snap-center cursor-pointer hover:scale-102 transition-transform duration-200"
                                    loading="lazy"
                                    onClick={() => openImageModal(image)}
                                    variants={imageVariants}
                                    initial="hidden"
                                    animate="visible"
                                    transition={{ delay: index * 0.1 }}
                                />
                            ))}
                        </div>
                        {images.length > 1 && (
                            <>
                                <motion.button
                                    onClick={() => handleScroll('left')}
                                    disabled={!canScrollLeft}
                                    className={`absolute left-2 sm:left-4 top-1/2 transform -translate-y-1/2 bg-black text-white p-2 sm:p-3 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity ${
                                        !canScrollLeft ? 'opacity-50 cursor-not-allowed' : ''
                                    }`}
                                    aria-label="Previous image"
                                    whileHover={{ scale: 1.1 }}
                                    whileTap={{ scale: 0.9 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    <ChevronLeft className="h-6 w-6 sm:h-8 sm:w-8" />
                                </motion.button>
                                <motion.button
                                    onClick={() => handleScroll('right')}
                                    disabled={!canScrollRight}
                                    className={`absolute right-2 sm:right-4 top-1/2 transform -translate-y-1/2 bg-black text-white p-2 sm:p-3 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity ${
                                        !canScrollRight ? 'opacity-50 cursor-not-allowed' : ''
                                    }`}
                                    aria-label="Next image"
                                    whileHover={{ scale: 1.1 }}
                                    whileTap={{ scale: 0.9 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    <ChevronRight className="h-6 w-6 sm:h-8 sm:w-8" />
                                </motion.button>
                            </>
                        )}
                    </div>

                    {selectedImage && (
                        <motion.div
                            className="fixed inset-0 bg-black/80 flex items-center justify-center z-50"
                            variants={modalVariants}
                            initial="hidden"
                            animate="visible"
                        >
                            <div ref={modalRef} className="relative max-w-5xl w-full p-2 sm:p-4">
                                <motion.button
                                    onClick={() => setSelectedImage(null)}
                                    className="absolute top-2 right-2 bg-black text-white p-2 rounded-full hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                    aria-label="Close image modal"
                                    whileHover={{ scale: 1.1 }}
                                    whileTap={{ scale: 0.9 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    <X className="h-6 w-6" />
                                </motion.button>
                                <img
                                    src={selectedImage}
                                    alt="Full-size property image"
                                    className="w-full max-h-[90vh] object-contain"
                                />
                            </div>
                        </motion.div>
                    )}

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 sm:gap-6">
                        <motion.div
                            className="md:col-span-2 space-y-4"
                            variants={detailsVariants}
                            initial="hidden"
                            animate="visible"
                        >
                            <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-4">
                                <h2 className="text-lg sm:text-xl font-bold text-black">{listing.title}</h2>
                                <span className="text-base sm:text-lg font-medium text-black mt-2 sm:mt-0">4.5 ★ (100)</span>
                            </div>
                            <p className="text-base sm:text-lg text-gray-500 mb-4">{listing.description}</p>
                            <p
                                className={`text-base sm:text-lg font-normal mb-4 ${
                                    listing.is_available ? 'text-green-600' : 'text-red-600'
                                }`}
                            >
                                {listing.is_available ? 'Available' : 'Booked'}
                            </p>
                            <p className="text-lg sm:text-xl font-bold text-black mb-4">
                                ₦{listing.price.toLocaleString()}/night
                            </p>
                            <p className="text-base sm:text-lg text-gray-600 mb-4">
                                Location: {listing.location}
                            </p>
                        </motion.div>
                        <motion.div
                            className="md:col-span-1"
                            variants={formVariants}
                            initial="hidden"
                            animate="visible"
                        >
                            <form onSubmit={handleBooking} className="space-y-4 bg-gray-50 p-4 sm:p-6 rounded-lg">
                                <h3 className="text-lg sm:text-xl font-bold text-black mb-4">Book Your Stay</h3>
                                <div>
                                    <label htmlFor="email" className="block text-sm sm:text-base font-medium text-black">
                                        Email
                                    </label>
                                    <input
                                        id="email"
                                        name="email"
                                        type="email"
                                        value={formData.email}
                                        onChange={handleInputChange}
                                        placeholder="Enter your email"
                                        className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-101 transition-all duration-150 text-base sm:text-lg"
                                        aria-label="Email address"
                                        required
                                        autoFocus={false}
                                    />
                                </div>
                                <div>
                                    <label htmlFor="name" className="block text-sm sm:text-base font-medium text-black">
                                        Name
                                    </label>
                                    <input
                                        id="name"
                                        name="name"
                                        type="text"
                                        value={formData.name}
                                        onChange={handleInputChange}
                                        placeholder="Enter your name"
                                        className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-101 transition-all duration-150 text-base sm:text-lg"
                                        aria-label="Full name"
                                        required
                                        autoFocus={false}
                                    />
                                </div>
                                <div>
                                    <label htmlFor="phone" className="block text-sm sm:text-base font-medium text-black">
                                        Phone
                                    </label>
                                    <input
                                        id="phone"
                                        name="phone"
                                        type="tel"
                                        value={formData.phone}
                                        onChange={handleInputChange}
                                        placeholder="Enter your phone number (e.g., +2341234567890)"
                                        className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-101 transition-all duration-150 text-base sm:text-lg"
                                        aria-label="Phone number"
                                        required
                                        autoFocus={false}
                                    />
                                </div>
                                <div>
                                    <label htmlFor="check-in" className="block text-sm sm:text-base font-medium text-black">
                                        Check-In
                                    </label>
                                    <DatePicker
                                        id="check-in"
                                        selected={checkIn}
                                        onChange={(date) => setCheckIn(date)}
                                        minDate={new Date()}
                                        placeholderText="Select date"
                                        className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-101 transition-all duration-150 text-base sm:text-lg"
                                        aria-label="Check-in date"
                                        wrapperClassName="w-full"
                                        popperClassName="z-50"
                                    />
                                </div>
                                <div>
                                    <label htmlFor="check-out" className="block text-sm sm:text-base font-medium text-black">
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
                                        className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-101 transition-all duration-150 text-base sm:text-lg"
                                        aria-label="Check-out date"
                                        wrapperClassName="w-full"
                                        popperClassName="z-50"
                                    />
                                </div>
                                {bookingError && (
                                    <motion.p
                                        className="text-red-500 text-sm sm:text-base"
                                        variants={messageVariants}
                                        initial="hidden"
                                        animate="visible"
                                    >
                                        {bookingError}
                                    </motion.p>
                                )}
                                {paymentStatus === 'success' && (
                                    <motion.p
                                        className="text-green-600 text-sm sm:text-base"
                                        variants={messageVariants}
                                        initial="hidden"
                                        animate="visible"
                                    >
                                        Booking created successfully!
                                    </motion.p>
                                )}
                                {paymentStatus === 'failed' && (
                                    <motion.p
                                        className="text-red-500 text-sm sm:text-base"
                                        variants={messageVariants}
                                        initial="hidden"
                                        animate="visible"
                                    >
                                        Payment failed. Please try again.
                                    </motion.p>
                                )}
                                <motion.div
                                    whileHover={{ scale: 1.05 }}
                                    whileTap={{ scale: 0.95 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    <PaystackButton
                                        className="w-full py-3 px-4 bg-black text-white font-medium rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 text-base sm:text-lg"
                                        {...paystackProps}
                                    />
                                </motion.div>
                            </form>
                        </motion.div>
                    </div>
                </div>
            </div>
        </motion.section>
    );
};

export default PropertyDetails;