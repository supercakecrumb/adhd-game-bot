import React from 'react';

interface BadgeProps {
  variant: 'difficulty' | 'mode';
  type: string;
}

const Badge: React.FC<BadgeProps> = ({ variant, type }) => {
  const getBadgeClass = () => {
    if (variant === 'difficulty') {
      switch (type) {
        case 'easy':
          return 'badge-difficulty-easy';
        case 'medium':
          return 'badge-difficulty-medium';
        case 'hard':
          return 'badge-difficulty-hard';
        default:
          return 'badge';
      }
    } else if (variant === 'mode') {
      switch (type) {
        case 'BINARY':
          return 'badge-mode-binary';
        case 'PARTIAL':
          return 'badge-mode-partial';
        case 'PER_MINUTE':
          return 'badge-mode-per-minute';
        default:
          return 'badge';
      }
    }
    return 'badge';
  };

  const getLabel = () => {
    if (variant === 'difficulty') {
      return type.charAt(0).toUpperCase() + type.slice(1);
    } else if (variant === 'mode') {
      switch (type) {
        case 'BINARY':
          return 'Binary';
        case 'PARTIAL':
          return 'Partial';
        case 'PER_MINUTE':
          return 'Per Minute';
        default:
          return type;
      }
    }
    return type;
  };

  return (
    <span className={getBadgeClass()}>
      {getLabel()}
    </span>
  );
};

export default Badge;