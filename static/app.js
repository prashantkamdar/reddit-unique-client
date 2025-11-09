// DOM elements
const usernameInput = document.getElementById('usernameInput');
const searchBtn = document.getElementById('searchBtn');
const loading = document.getElementById('loading');
const error = document.getElementById('error');
const resultsInfo = document.getElementById('resultsInfo');
const postsContainer = document.getElementById('postsContainer');

// Event listeners
searchBtn.addEventListener('click', searchUser);
usernameInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        searchUser();
    }
});

async function searchUser() {
    const username = usernameInput.value.trim();
    
    if (!username) {
        showError('Please enter a Reddit username');
        return;
    }
    
    // Clear previous results
    hideError();
    hideResults();
    showLoading();
    
    try {
        const response = await fetch(`/api/user/${encodeURIComponent(username)}`);
        const data = await response.json();
        
        hideLoading();
        
        if (data.error) {
            showError(data.error);
            return;
        }
        
        if (!data.posts || data.posts.length === 0) {
            showError(`No posts found for user: ${username}`);
            return;
        }
        
        displayPosts(data.posts, username);
    } catch (err) {
        hideLoading();
        showError('Failed to fetch posts. Please try again.');
        console.error('Error:', err);
    }
}

function displayPosts(posts, username) {
    resultsInfo.textContent = `Found ${posts.length} unique post${posts.length !== 1 ? 's' : ''} by u/${username}`;
    resultsInfo.classList.remove('hidden');
    
    postsContainer.innerHTML = '';
    
    posts.forEach(post => {
        const postCard = createPostCard(post);
        postsContainer.appendChild(postCard);
    });
}

function createPostCard(post) {
    const card = document.createElement('a');
    card.className = 'post-card';
    card.href = `https://old.reddit.com${post.permalink}`;
    card.target = '_blank';
    card.rel = 'noopener noreferrer';
    
    // Format time
    const timeAgo = formatTimeAgo(post.created_utc);
    
    // Format numbers
    const score = formatNumber(post.score);
    const comments = formatNumber(post.num_comments);
    
    // Check if we should show thumbnail
    const showThumbnail = post.thumbnail && 
                          post.thumbnail !== 'self' && 
                          post.thumbnail !== 'default' && 
                          post.thumbnail !== 'nsfw' &&
                          post.thumbnail.startsWith('http');
    
    // Truncate selftext
    const selftext = post.selftext ? truncateText(post.selftext, 150) : '';
    
    card.innerHTML = `
        <div class="post-header">
            <span class="subreddit-name">r/${post.subreddit}</span>
            <span class="separator">•</span>
            <span class="post-time">${timeAgo}</span>
        </div>
        <h2 class="post-title">${escapeHtml(post.title)}</h2>
        <div class="post-content">
            ${showThumbnail ? `<img src="${post.thumbnail}" alt="Post thumbnail" class="post-thumbnail">` : ''}
            ${selftext ? `<div class="post-text"><p class="post-selftext">${escapeHtml(selftext)}</p></div>` : ''}
        </div>
        <div class="post-meta">
            <div class="meta-item upvotes">
                <svg class="meta-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 4l-8 8h5v8h6v-8h5z"/>
                </svg>
                ${score}
            </div>
            <div class="meta-item">
                <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                </svg>
                ${comments}
            </div>
        </div>
    `;
    
    return card;
}

function formatTimeAgo(timestamp) {
    const seconds = Math.floor(Date.now() / 1000 - timestamp);
    
    if (seconds < 60) return 'just now';
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
    if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;
    if (seconds < 2592000) return `${Math.floor(seconds / 604800)}w ago`;
    if (seconds < 31536000) return `${Math.floor(seconds / 2592000)}mo ago`;
    return `${Math.floor(seconds / 31536000)}y ago`;
}

function formatNumber(num) {
    if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
    return num.toString();
}

function truncateText(text, maxLength) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showLoading() {
    loading.classList.remove('hidden');
}

function hideLoading() {
    loading.classList.add('hidden');
}

function showError(message) {
    error.textContent = message;
    error.classList.remove('hidden');
}

function hideError() {
    error.classList.add('hidden');
}

function hideResults() {
    resultsInfo.classList.add('hidden');
    postsContainer.innerHTML = '';
}

