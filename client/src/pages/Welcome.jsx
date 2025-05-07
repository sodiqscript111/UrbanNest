import React from 'react';
import Nav from "../component/Nav.jsx";
import Header from "../component/Header.jsx";
import Features from "../component/Features.jsx";
import GuestFavorites  from "../component/Guestfavorite.jsx";
import BudgetFriendly  from "../component/Budgetfriendly.jsx";
import FAQ from "../component/Faqs.jsx";

const Welcome = () => {
    return (
        <>
        <Nav />
            <Header />
            <GuestFavorites />
            <Features/>
            <BudgetFriendly/>
            <FAQ/>
        </>
    )
}
export default Welcome;