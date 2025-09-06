# ADHD Game Bot UI

A React TypeScript frontend for the ADHD Game Bot application.

## Features

- Dashboard with today's quests
- Quest catalog with filtering
- Shop for purchasing items
- User profile
- Admin panel for creating quests and shop items
- Responsive design for mobile and desktop
- Dark theme UI with Tailwind CSS

## Tech Stack

- React 18 with TypeScript
- Vite for fast development
- Tailwind CSS for styling
- React Router for navigation
- Lucide React for icons

## Getting Started

### Prerequisites

- Node.js (v16 or higher)
- npm or yarn

### Installation

1. Clone the repository
2. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
3. Install dependencies:
   ```bash
   npm install
   ```
   or
   ```bash
   yarn install
   ```

### Development

To start the development server:

```bash
npm run dev
```
or
```bash
yarn dev
```

The application will be available at http://localhost:3000

### Building for Production

To create a production build:

```bash
npm run build
```
or
```bash
yarn build
```

### Environment Variables

Create a `.env` file in the frontend directory with the following variables:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## Project Structure

```
src/
├── components/     # Reusable UI components
├── pages/          # Page components
├── services/       # API services and auth
├── types/          # TypeScript types
├── utils/          # Utility functions
├── App.tsx         # Main app component
└── main.tsx        # Entry point
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Design System

The UI uses a dark theme with the following color palette:

- Background: slate-900
- Surface: slate-800
- Border: slate-700
- Text: slate-100
- Muted: slate-400
- Primary: violet-600
- Success: green-500
- Warning: amber-500
- Danger: rose-500

## Responsive Design

The application is fully responsive:
- Mobile: Bottom navigation bar
- Desktop: Left sidebar navigation
- Adaptive layouts for all screen sizes

## Accessibility

- Proper contrast ratios for readability
- Keyboard navigation support
- Semantic HTML structure
- ARIA attributes where needed