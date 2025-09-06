import React, { useState } from 'react';

interface ShopItem {
  id: number;
  name: string;
  price: number;
  stock?: number;
  category: string;
}

const Shop: React.FC = () => {
  const [selectedItem, setSelectedItem] = useState<ShopItem | null>(null);
  const [quantity, setQuantity] = useState(1);

  // Mock shop items
  const shopItems: ShopItem[] = [
    {
      id: 1,
      name: 'Extra Life',
      price: 200,
      stock: 5,
      category: 'Power-ups',
    },
    {
      id: 2,
      name: 'Focus Boost',
      price: 150,
      category: 'Power-ups',
    },
    {
      id: 3,
      name: 'Time Extension',
      price: 300,
      stock: 3,
      category: 'Power-ups',
    },
    {
      id: 4,
      name: 'Streak Protector',
      price: 500,
      stock: 2,
      category: 'Power-ups',
    },
    {
      id: 5,
      name: 'Double Points',
      price: 1000,
      stock: 1,
      category: 'Power-ups',
    },
    {
      id: 6,
      name: 'Custom Theme',
      price: 750,
      category: 'Cosmetics',
    },
  ];

  const handlePurchase = (item: ShopItem) => {
    setSelectedItem(item);
    setQuantity(1);
  };

  const confirmPurchase = () => {
    if (selectedItem) {
      console.log(`Purchasing ${quantity}x ${selectedItem.name}`);
      // In a real app, we would call the API here
      setSelectedItem(null);
    }
  };

  return (
    <div className="min-h-screen bg-slate-900">
      {/* Header */}
      <div className="sticky top-0 z-10 border-b border-slate-700 bg-slate-900/80 backdrop-blur">
        <div className="max-w-4xl mx-auto px-4 py-3">
          <h1 className="text-2xl font-bold">Shop</h1>
        </div>
      </div>

      {/* Shop Content */}
      <div className="max-w-4xl mx-auto px-4 py-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {shopItems.map((item) => (
            <div key={item.id} className="card">
              <h3 className="font-bold text-lg mb-2">{item.name}</h3>
              <p className="text-slate-400 text-sm mb-4">{item.category}</p>
              <div className="flex items-center justify-between">
                <span className="text-violet-400 font-medium">{item.price} points</span>
                {item.stock !== undefined && (
                  <span className="text-sm text-slate-400">
                    {item.stock} left
                  </span>
                )}
              </div>
              <button
                onClick={() => handlePurchase(item)}
                className="btn-primary w-full mt-4"
                disabled={item.stock === 0}
              >
                {item.stock === 0 ? 'Sold Out' : 'Purchase'}
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Purchase Dialog */}
      {selectedItem && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="card w-full max-w-md">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold">Purchase {selectedItem.name}</h3>
              <button 
                onClick={() => setSelectedItem(null)}
                className="text-slate-400 hover:text-slate-100"
              >
                <svg 
                  xmlns="http://www.w3.org/2000/svg" 
                  width="20" 
                  height="20" 
                  viewBox="0 0 24 24" 
                  fill="none" 
                  stroke="currentColor" 
                  strokeWidth="2" 
                  strokeLinecap="round" 
                  strokeLinejoin="round"
                >
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">
                  Quantity
                </label>
                <input
                  type="number"
                  min="1"
                  max={selectedItem.stock || 99}
                  value={quantity}
                  onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                  className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                  disabled={selectedItem.stock === 0}
                />
                {selectedItem.stock !== undefined && (
                  <p className="text-xs text-slate-400 mt-1">
                    Available: {selectedItem.stock}
                  </p>
                )}
              </div>
              
              <div className="flex justify-between items-center pt-4 border-t border-slate-700">
                <span className="font-medium">Total:</span>
                <span className="text-xl font-bold text-violet-400">
                  {quantity * selectedItem.price} points
                </span>
              </div>
            </div>
            
            <div className="flex justify-end space-x-3 mt-6">
              <button
                onClick={() => setSelectedItem(null)}
                className="btn-secondary"
              >
                Cancel
              </button>
              <button
                onClick={confirmPurchase}
                className="btn-primary"
                disabled={selectedItem.stock === 0}
              >
                Confirm Purchase
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Shop;