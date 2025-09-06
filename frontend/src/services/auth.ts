const USER_ID_KEY = 'adhd-game-user-id';

export const getUserID = (): string | null => {
  return localStorage.getItem(USER_ID_KEY);
};

export const setUserID = (userId: string): void => {
  localStorage.setItem(USER_ID_KEY, userId);
};

export const clearUserID = (): void => {
  localStorage.removeItem(USER_ID_KEY);
};

export const getAuthHeaders = (): Record<string, string> => {
  const userId = getUserID();
  if (userId) {
    return {
      'X-User-ID': userId,
    };
  }
  return {};
};