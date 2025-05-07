import { useState } from 'react';
import { ChevronDown, ChevronUp } from 'lucide-react';

export default function FAQ() {
    const [openIndex, setOpenIndex] = useState(null);

    const faqs = [
        {
            question: 'How do I sign up for UrbanNest?',
            answer:
                'To sign up, visit the signup page, enter your email, password, first name, last name, and phone number. Your email will be verified in real-time using ZeroBounce. Ensure your email is valid and not disposable. Once verified, submit the form to create your account.',
        },
        {
            question: 'Why am I getting an "Email does not exist or is invalid" error during signup?',
            answer:
                'This error occurs if our email verification service (ZeroBounce) determines your email is invalid or doesn’t exist. Ensure your email is typed correctly. If the issue persists, it might be due to temporary service issues or restrictions on certain email providers. Contact support at support@urbannest.com.',
        },
        {
            question: 'How do I list a property on UrbanNest?',
            answer:
                'After signing up as a lister, log in and navigate to “My Listings.” Click “Add Listing,” provide details like title, description, price, and upload images using the provided upload URL. Your listing will be live once submitted.',
        },
        {
            question: 'How does the booking process work?',
            answer:
                'Browse listings, select a property, and choose your dates. Complete the payment via Paystack. Once the payment is verified, your booking is confirmed, and you’ll receive a confirmation email. You can view your bookings under “My Bookings.”',
        },
        {
            question: 'What payment methods are accepted?',
            answer:
                'We use Paystack for secure payments, supporting major credit/debit cards and bank transfers (in supported regions). Ensure your payment matches the expected amount for the booking duration.',
        },
        {
            question: 'Can I edit or delete my listing?',
            answer:
                'Yes, log in as a lister, go to “My Listings,” select the listing, and choose “Edit” to update details or “Delete” to remove it. Changes are reflected immediately.',
        },
        {
            question: 'What should I do if I encounter an error during signup or booking?',
            answer:
                'Check your internet connection and ensure all fields are filled correctly. For signup issues, verify your email format. For booking issues, confirm your payment details. If problems persist, contact support at support@urbannest.com with error details.',
        },
    ];

    const toggleFAQ = (index) => {
        setOpenIndex(openIndex === index ? null : index);
    };

    return (
        <div className="min-h-screen bg-gray-100 flex items-start py-12">
            <div className="max-w-4xl w-full mx-auto px-6">
                <h2 className="text-3xl font-bold text-gray-800 text-center mb-8">Frequently Asked Questions</h2>
                <div className="space-y-4">
                    {faqs.map((faq, index) => (
                        <div key={index} className="bg-white rounded-md shadow-sm overflow-hidden">
                            <button
                                className="w-full flex justify-between items-center p-4 text-left transition-colors duration-200 ease-in-out hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                                onClick={() => toggleFAQ(index)}
                            >
                                <span className="text-lg font-medium text-gray-800">{faq.question}</span>
                                {openIndex === index ? (
                                    <ChevronUp className="w-5 h-5 text-indigo-500" />
                                ) : (
                                    <ChevronDown className="w-5 h-5 text-gray-500" />
                                )}
                            </button>
                            {openIndex === index && (
                                <div className="p-4 bg-gray-50">
                                    <p className="text-gray-700">{faq.answer}</p>
                                </div>
                            )}
                        </div>
                    ))}
                </div>
                <p className="text-center text-gray-600 mt-8 text-lg">
                    Still have questions?{' '}
                    <a href="mailto:support@urbannest.com" className="text-indigo-500 font-medium hover:underline">
                        Contact Support
                    </a>
                </p>
            </div>
        </div>
    );
}