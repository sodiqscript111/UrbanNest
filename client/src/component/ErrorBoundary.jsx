import React, { Component } from 'react';
import { Link } from 'react-router-dom';

class ErrorBoundary extends Component {
    state = { hasError: false, error: null };

    static getDerivedStateFromError(error) {
        return { hasError: true, error };
    }

    componentDidCatch(error, errorInfo) {
        console.error('ErrorBoundary caught an error:', error, errorInfo);
    }

    render() {
        if (this.state.hasError) {
            return (
                <div className="p-4 max-w-7xl mx-auto font-sans">
                    <h1 className="text-2xl font-bold text-red-500 mb-4">Something went wrong</h1>
                    <p className="text-gray-600 mb-4">
                        An unexpected error occurred. Please try again or return to the homepage.
                    </p>
                    <Link
                        to="/"
                        className="inline-block px-4 py-2 bg-black text-white font-medium rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                    >
                        Back to Home
                    </Link>
                </div>
            );
        }

        return this.props.children;
    }
}

export default ErrorBoundary;