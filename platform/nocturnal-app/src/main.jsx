import { createRoot } from 'react-dom/client'
import './index.scss'
import App from './App.jsx'
import 'materialize-css/dist/css/materialize.min.css';
import 'materialize-css/dist/js/materialize.min.js';

createRoot(document.getElementById('root')).render(<App />)
