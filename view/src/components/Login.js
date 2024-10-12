import React, { useState } from 'react';
import axios from 'axios';

const Login = () => {
    const [credentials, setCredentials] = useState({
        phone: '',
        password: '',
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setCredentials({
            ...credentials,
            [name]: value,
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8080/login', credentials);
            const token = response.data.token;
            localStorage.setItem('jwt', token); // Сохраняем токен
            console.log('Logged in successfully:', response.data);
        } catch (error) {
            console.error('Error logging in:', error);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input name="phone" type="text" placeholder="Phone" onChange={handleChange} required />
            <input name="password" type="password" placeholder="Password" onChange={handleChange} required />
            <button type="submit">Login</button>
        </form>
    );
};

export default Login;
