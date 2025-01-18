import React, { useState } from 'react';
import axios from 'axios';
import './ShortenUrl.css';
import heroImage from '../assets/share-link.svg'; // Ensure this path is correct
import linkIcon from '../assets/link.png'; // Ensure this path is correct
import copyIcon from '../assets/copy-icon.svg'; // Ensure this path is correct

const ShortenUrl = () => {
  const [url, setUrl] = useState<string>('');
  const [shortUrl, setShortUrl] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [copied, setCopied] = useState<boolean>(false);

  // Validate URL
  const isValidUrl = (url: string): boolean => {
    try {
      new URL(url);
      return true;
    } catch (err) {
      return false;
    }
  };

  const handleShorten = async (e?: React.FormEvent): Promise<void> => {
    if (e) {
      e.preventDefault(); // Prevent form submission from refreshing the page
    }

    if (!url) {
      setError('Please enter a valid URL');
      return;
    }

    if (!isValidUrl(url)) {
      setError('Please enter a valid URL starting with http:// or https://');
      setUrl(''); // Clear the input box
      return;
    }

    setLoading(true);
    setError('');
    setCopied(false);

    try {
      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080'; // Fallback to localhost
      console.log("Api used: ",apiUrl);
      const response = await axios.post<{ short_url: string }>(`${apiUrl}/shorten`, { url });
      setShortUrl(response.data.short_url);
    } catch (err) {
      console.error('Error:', err);
      setError('Failed to shorten the URL. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = (): void => {
    if (shortUrl) {
      navigator.clipboard.writeText(shortUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000); // Reset copied state after 2 seconds
    }
  };

  return (
    <div className="container">
      <div className="main-content">
        <div className="hero-section">
          <img src={heroImage} alt="Shortly" className="hero-image" />
          <h1>Welcome to Shortly!</h1>
          <p className="subtitle">Shorten your links in seconds!</p>
        </div>
        <div className="card">
          <form onSubmit={handleShorten}>
            <div className="input-group">
              <img src={linkIcon} alt="Link Icon" className="link-icon" />
              <input
                type="text"
                value={url}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setUrl(e.target.value)}
                placeholder="Paste your long URL here..."
                disabled={loading}
                aria-label="Enter URL to shorten"
              />
              <button type="submit" disabled={loading} aria-label="Shorten URL">
                {loading ? <div className="spinner"></div> : 'Shortly!'}
              </button>
            </div>
          </form>
          {error && (
            <div className="alert">
              {error}
              <button onClick={() => setError('')} className="close-alert" aria-label="Close alert">
                &times;
              </button>
            </div>
          )}
          {shortUrl && (
            <div className="result">
              <p>Here's your shortened URL:</p>
              <div className="shortened-url-container">
                <a href={shortUrl} target="_blank" rel="noopener noreferrer" aria-label="Shortened URL">
                  {shortUrl}
                </a>
                <button onClick={handleCopy} className="copy-button" aria-label="Copy shortened URL">
                  <img src={copyIcon} alt="Copy" className="copy-icon" />
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
      <footer className="footer">
        <p>
          Made with ❤️ by{' '}
          <a href="https://github.com/yash6318" target="_blank" rel="noopener noreferrer">
            Yash Agrawal
          </a>
        </p>
        <p>© 2025 Shortly. All rights reserved.</p>
      </footer>
    </div>
  );
};

export default ShortenUrl;
