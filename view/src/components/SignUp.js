import React, { useState } from 'react';
import axios from 'axios';

const SignUp = () => {
    const [formData, setFormData] = useState({
        phone: '',
        email: '',
        name: '',
        surname: '',
        password: '',
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value,
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8080/sign-up', formData);
            console.log(response.data);
        } catch (error) {
            console.error('Error signing up:', error);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input name="phone" type="text" placeholder="Phone" onChange={handleChange} required />
            <input name="email" type="email" placeholder="Email" onChange={handleChange} required />
            <input name="name" type="text" placeholder="Name" onChange={handleChange} required />
            <input name="surname" type="text" placeholder="Surname" onChange={handleChange} required />
            <input name="password" type="password" placeholder="Password" onChange={handleChange} required />
            <button type="submit">Sign Up</button>
        </form>
    );
};

export default SignUp;
