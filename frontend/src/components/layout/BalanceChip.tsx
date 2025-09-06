import React from 'react';

interface BalanceChipProps {
  balance: string;
}

const BalanceChip: React.FC<BalanceChipProps> = ({ balance }) => {
  // Format the balance with commas
  const formattedBalance = new Intl.NumberFormat().format(parseInt(balance) || 0);

  return (
    <div className="flex items-center space-x-1 rounded-full bg-slate-800 px-3 py-1 text-sm font-medium">
      <span>💰</span>
      <span>{formattedBalance}</span>
    </div>
  );
};

export default BalanceChip;