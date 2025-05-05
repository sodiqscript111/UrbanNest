import { useState } from 'react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';

const Header = () => {
    const [location, setLocation] = useState('');
    const [checkIn, setCheckIn] = useState(null);
    const [checkOut, setCheckOut] = useState(null);
    const [error, setError] = useState('');

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

    return (
        <header className="bg-white font-sans">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                <div className="flex flex-col md:grid md:grid-cols-[1fr_3fr] gap-8">
                    {/* Image */}
                    <div className="order-1 md:order-2 relative h-[26.375rem] md:h-[35.2rem]">
                        <img
                            src="https://i.imgur.com/wD82gGh.jpeg"
                            alt="Modern apartment in Lagos"
                            className="w-full h-full object-cover rounded-lg"
                            loading="lazy"
                        />
                    </div>

                    {/* Card (Search Form) */}
                    <div className="order-2 md:order-1 relative z-10 md:-mr-12">
                        <div className="bg-white shadow-lg rounded-2xl p-14 max-w-3xl mt-8 scale-105">
                        <h1 className="text-3xl font-bold text-black mb-6">
                                Find Your Perfect Nest
                            </h1>
                            <form onSubmit={handleSearch} className="space-y-4">
                                <div>
                                    <label htmlFor="location" className="block text-sm font-medium text-black">
                                        Location
                                    </label>
                                    <input
                                        type="text"
                                        id="location"
                                        value={location}
                                        onChange={(e) => setLocation(e.target.value)}
                                        placeholder="e.g., Lagos, Nigeria"
                                        className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                        aria-describedby="location-description"
                                    />
                                    <p id="location-description" className="text-xs text-gray-500 mt-1">
                                        Enter a city or neighborhood in Lagos.
                                    </p>
                                </div>
                                <div className="flex space-x-4">
                                    <div className="flex-1">
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
                                    <div className="flex-1">
                                        <label htmlFor="check-out" className="block text-sm font-medium text-black">
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
                                            className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md text-black placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                            aria-label="Check-out date"
                                        />
                                    </div>
                                </div>
                                {error && <p className="text-red-500 text-sm">{error}</p>}
                                <button
                                    type="submit"
                                    className="w-full py-2 px-4 bg-black text-white font-medium rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                >
                                    Search
                                </button>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </header>
    );
};

export default Header;