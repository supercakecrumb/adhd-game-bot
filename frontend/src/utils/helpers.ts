export const generateId = (): string => {
  return crypto.randomUUID();
};

export const formatPoints = (points: string): string => {
  try {
    const num = parseFloat(points);
    return num.toLocaleString();
  } catch {
    return points;
  }
};

export const calculateReward = (quest: any, value?: number): string => {
  switch (quest.mode) {
    case 'BINARY':
      return quest.points_award;
    case 'PARTIAL':
      if (value !== undefined) {
        const maxPoints = parseFloat(quest.points_award);
        return (maxPoints * (value / 100)).toFixed(2);
      }
      return `Up to ${quest.points_award}`;
    case 'PER_MINUTE':
      if (value !== undefined && quest.rate_points_per_min) {
        const rate = parseFloat(quest.rate_points_per_min);
        const minutes = value;
        let points = rate * minutes;
        
        // Apply cap if exists
        if (quest.max_minutes) {
          const maxPoints = rate * quest.max_minutes;
          points = Math.min(points, maxPoints);
        }
        
        return points.toFixed(2);
      }
      if (quest.rate_points_per_min) {
        return `${quest.rate_points_per_min} per minute`;
      }
      return '0';
    default:
      return '0';
  }
};