import { useState } from 'react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';

const Header = () => {
    const [location, setLocation] = useState('');
    const [checkIn, setCheckIn] = useState(null);
    const [checkOut, setCheckOut] = useState(null);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleFindNest = () => {
        navigate('/home');
    };

    const handleSearch = (e) => {
        e.preventDefault();
        if (!location || !checkIn || !checkOut) {
            setError('Please fill all fields');
            return;
        }
        setError('');
        const query = {
            location,
            checkIn: checkIn ? checkIn.toISOString().split('T')[0] : '',
            checkOut: checkOut ? checkOut.toISOString().split('T')[0] : '',
        };
        console.log('Search query:', query);
    };

    // Animation variants for form card
    const cardVariants = {
        hidden: { opacity: 0, y: 20 },
        visible: { opacity: 1, y: 0, transition: { duration: 0.6, ease: 'easeOut' } },
    };

    // Animation variants for image
    const imageVariants = {
        hidden: { opacity: 0 },
        visible: { opacity: 1, transition: { duration: 0.8, ease: 'easeIn' } },
    };

    // Animation variants for error message
    const errorVariants = {
        hidden: { opacity: 0 },
        visible: { opacity: 1, transition: { duration: 0.3 } },
    };

    return (
        <header className="bg-white font-sans">
            <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 sm:py-10 lg:px-8 lg:py-12">
                <div className="flex flex-col md:grid md:grid-cols-[1fr_3fr] gap-6 md:gap-8">
                    {/* Image */}
                    <motion.div
                        className="order-1 md:order-2 relative h-64 md:h-[35.2rem]"
                        variants={imageVariants}
                        initial="hidden"
                        animate="visible"
                    >
                        <img
                            src="https://i.imgur.com/wD82gGh.jpeg"
                            alt="Modern apartment in Lagos"
                            className="w-full h-full object-cover rounded-lg"
                            loading="lazy"
                        />
                    </motion.div>

                    {/* Card (Search Form) */}
                    <div className="order-2 md:order-1 relative z-10 md:-mr-12">
                        <motion.div
                            className="bg-white shadow-lg rounded-2xl p-6 md:p-14 max-w-full md:max-w-3xl mt-6 md:mt-8 md:scale-105"
                            variants={cardVariants}
                            initial="hidden"
                            animate="visible"
                        >
                            <h1 className="text-2xl md:text-3xl font-bold text-black mb-4 md:mb-6">
                                Find Your Perfect Nest
                            </h1>
                            <form onSubmit={handleSearch} className="space-y-4">
                                <div>
                                    <label htmlFor="location" className="block text-xs md:text-sm font-medium text-black">
                                        Location
                                    </label>
                                    <input
                                        type="text"
                                        id="location"
                                        value={location}
                                        onChange={(e) => setLocation(e.target.value)}
                                        placeholder="e.g., Lagos, Nigeria"
                                        className="mt-1 w-full px-3 py-1.5 md:py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-102 transition-all duration-200"
                                        aria-describedby="location-description"
                                    />
                                    <p id="location-description" className="text-[0.65rem] md:text-xs text-gray-500 mt-1">
                                        Enter a city or neighborhood in Lagos.
                                    </p>
                                </div>
                                <div className="flex flex-col md:flex-row md:space-x-4 space-y-4 md:space-y-0">
                                    <div className="flex-1">
                                        <label htmlFor="check-in" className="block text-xs md:text-sm font-medium text-black">
                                            Check-In
                                        </label>
                                        <DatePicker
                                            id="check-in"
                                            selected={checkIn}
                                            onChange={(date) => setCheckIn(date)}
                                            minDate={new Date()}
                                            placeholderText="Select date"
                                            className="mt-1 w-full px-3 py-1.5 md:py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-102 transition-all duration-200"
                                            aria-label="Check-in date"
                                        />
                                    </div>
                                    <div className="flex-1">
                                        <label htmlFor="check-out" className="block text-xs md:text-sm font-medium text-black">
                                            Check-Out
                                        </label>
                                        <DatePicker
                                            id="check-out"
                                            selected={checkOut}
                                            onChange={(date) => {
                                                if (checkIn && date < checkIn) {
                                                    setError('Check-out must be after check-in');
                                                    return;
                                                }
                                                setCheckOut(date);
                                                setError('');
                                            }}
                                            minDate={checkIn || new Date()}
                                            placeholderText="Select date"
                                            className="mt-1 w-full px-3 py-1.5 md:py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 focus:scale-102 transition-all duration-200"
                                            aria-label="Check-out date"
                                        />
                                    </div>
                                </div>
                                {error && (
                                    <motion.p
                                        className="text-sm text-red-500"
                                        variants={errorVariants}
                                        initial="hidden"
                                        animate="visible"
                                    >
                                        {error}
                                    </motion.p>
                                )}
                                <motion.button
                                    type="submit"
                                    className="w-full py-2 px-4 bg-black text-white font-medium rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 min-h-[2.5rem]"
                                    whileHover={{ scale: 1.05 }}
                                    whileTap={{ scale: 0.95 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    Search
                                </motion.button>
                                <motion.button
                                    type="button"
                                    onClick={handleFindNest}
                                    className="w-full py-2 px-4 bg-black text-white font-medium rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500 mt-2 md:mt-4 min-h-[2.5rem]"
                                    whileHover={{ scale: 1.05 }}
                                    whileTap={{ scale: 0.95 }}
                                    transition={{ duration: 0.2 }}
                                    aria-label="Find your nest"
                                >
                                    Find Your Nest
                                </motion.button>
                            </form>
                        </motion.div>
                    </div>
                </div>
            </div>
        </header>
    );
};

export default Header;