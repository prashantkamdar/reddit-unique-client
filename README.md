# Reddit Client - Apollo Style

A clean, modern Reddit client built with Go that displays unique posts from any Reddit user. Features an Apollo-inspired UI with a beautiful, minimalist design.

## Features

- **Unique Posts Only**: Automatically deduplicates posts (if the same post appears in multiple subreddits, it's shown only once)
- **Chronological Order**: Posts are sorted from newest to oldest
- **Apollo-Inspired UI**: Clean, modern interface inspired by the beloved Apollo Reddit client
- **No Authentication Required**: Uses Reddit's public JSON endpoints
- **Fast & Lightweight**: Built with Go for optimal performance
- **Responsive Design**: Works great on desktop and mobile

## Installation

### Prerequisites
- Go 1.21 or higher

### Running the Application

1. Clone or navigate to this directory
2. Run the application:

```bash
go run main.go
```

3. Open your browser and visit: `http://localhost:8080`

## Usage

1. Enter any Reddit username in the search box
2. Click "Search" or press Enter
3. Browse through the user's unique posts in a beautiful, card-based layout
4. Click any post to open it on Reddit

## Technical Details

### Backend (Go)
- Uses Reddit's public JSON API (`/user/{username}/submitted.json`)
- Fetches multiple pages to get comprehensive post history
- Deduplicates posts by post ID
- Sorts posts chronologically
- Implements rate limiting to be respectful to Reddit's servers

### Frontend
- Pure HTML, CSS, and JavaScript (no frameworks needed)
- Apollo-inspired design with:
  - Clean card-based layout
  - Smooth animations and transitions
  - Blue accent colors (#007aff)
  - Proper spacing and typography
  - Responsive design

### File Structure
```
reddit-client/
├── main.go           # Go backend server
├── go.mod            # Go module definition
├── README.md         # This file
└── static/
    ├── index.html    # Main HTML page
    ├── style.css     # Apollo-inspired styles
    └── app.js        # Frontend JavaScript
```

## API Endpoint

The application exposes a single API endpoint:

- `GET /api/user/{username}` - Fetches all unique posts for a Reddit user

Response format:
```json
{
  "posts": [
    {
      "id": "abc123",
      "title": "Post title",
      "subreddit": "example",
      "author": "username",
      "score": 100,
      "num_comments": 10,
      "created_utc": 1699564800,
      "url": "https://...",
      "permalink": "/r/example/comments/...",
      "thumbnail": "https://...",
      "selftext": "Post content..."
    }
  ]
}
```

## License

MIT

