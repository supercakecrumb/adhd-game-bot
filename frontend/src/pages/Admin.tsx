import React, { useState } from 'react';

const Admin: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'quest' | 'shop'>('quest');
  
  // Quest form state
  const [questForm, setQuestForm] = useState({
    title: '',
    description: '',
    category: 'daily',
    difficulty: 'easy',
    mode: 'BINARY',
    points_award: '',
    rate_points_per_min: '',
    min_minutes: '',
    max_minutes: '',
    cooldown_sec: '3600',
    streak_enabled: true,
  });
  
  // Shop item form state
  const [shopForm, setShopForm] = useState({
    name: '',
    code: '',
    description: '',
    price: '',
    stock: '',
    is_active: true,
  });

  const handleQuestFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const checked = type === 'checkbox' ? (e.target as HTMLInputElement).checked : undefined;
    
    setQuestForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleShopFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const checked = type === 'checkbox' ? (e.target as HTMLInputElement).checked : undefined;
    
    setShopForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleQuestSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Creating quest:', questForm);
    // In a real app, we would call the API here
    alert('Quest created successfully!');
  };

  const handleShopSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Creating shop item:', shopForm);
    // In a real app, we would call the API here
    alert('Shop item created successfully!');
  };

  return (
    <div className="min-h-screen bg-slate-900">
      {/* Header */}
      <div className="sticky top-0 z-10 border-b border-slate-700 bg-slate-900/80 backdrop-blur">
        <div className="max-w-4xl mx-auto px-4 py-3">
          <h1 className="text-2xl font-bold">Admin Panel</h1>
        </div>
      </div>

      {/* Tabs */}
      <div className="max-w-4xl mx-auto px-4 py-6">
        <div className="flex space-x-4 border-b border-slate-700 mb-6">
          <button
            onClick={() => setActiveTab('quest')}
            className={`pb-3 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'quest'
                ? 'border-violet-500 text-violet-400'
                : 'border-transparent text-slate-400 hover:text-slate-300'
            }`}
          >
            Create Quest
          </button>
          <button
            onClick={() => setActiveTab('shop')}
            className={`pb-3 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'shop'
                ? 'border-violet-500 text-violet-400'
                : 'border-transparent text-slate-400 hover:text-slate-300'
            }`}
          >
            Create Shop Item
          </button>
        </div>

        {/* Create Quest Form */}
        {activeTab === 'quest' && (
          <div className="card">
            <h2 className="text-xl font-bold mb-4">Create New Quest</h2>
            <form onSubmit={handleQuestSubmit} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">Title</label>
                  <input
                    type="text"
                    name="title"
                    value={questForm.title}
                    onChange={handleQuestFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                    required
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">Category</label>
                  <select
                    name="category"
                    value={questForm.category}
                    onChange={handleQuestFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                  >
                    <option value="daily">Daily</option>
                    <option value="weekly">Weekly</option>
                    <option value="adhoc">One-time</option>
                  </select>
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">Difficulty</label>
                  <select
                    name="difficulty"
                    value={questForm.difficulty}
                    onChange={handleQuestFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                  >
                    <option value="easy">Easy</option>
                    <option value="medium">Medium</option>
                    <option value="hard">Hard</option>
                  </select>
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">Mode</label>
                  <select
                    name="mode"
                    value={questForm.mode}
                    onChange={handleQuestFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                  >
                    <option value="BINARY">Binary</option>
                    <option value="PARTIAL">Partial</option>
                    <option value="PER_MINUTE">Per Minute</option>
                  </select>
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-1">Description</label>
                <textarea
                  name="description"
                  value={questForm.description}
                  onChange={handleQuestFormChange}
                  rows={3}
                  className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                  required
                />
              </div>
              
              {/* Mode-specific fields */}
              {questForm.mode === 'BINARY' || questForm.mode === 'PARTIAL' ? (
                <div>
                  <label className="block text-sm font-medium mb-1">
                    {questForm.mode === 'BINARY' ? 'Points Award' : 'Max Points'}
                  </label>
                  <input
                    type="number"
                    name="points_award"
                    value={questForm.points_award}
                    onChange={handleQuestFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                    required
                  />
                </div>
              ) : null}
              
              {questForm.mode === 'PER_MINUTE' && (
                <>
                  <div>
                    <label className="block text-sm font-medium mb-1">Points per Minute</label>
                    <input
                      type="number"
                      name="rate_points_per_min"
                      value={questForm.rate_points_per_min}
                      onChange={handleQuestFormChange}
                      step="0.1"
                      className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                      required
                    />
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium mb-1">Min Minutes (optional)</label>
                      <input
                        type="number"
                        name="min_minutes"
                        value={questForm.min_minutes}
                        onChange={handleQuestFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                      />
                    </div>
                    
                    <div>
                      <label className="block text-sm font-medium mb-1">Max Minutes (optional)</label>
                      <input
                        type="number"
                        name="max_minutes"
                        value={questForm.max_minutes}
                        onChange={handleQuestFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                      />
                    </div>
                  </div>
                </>
              )}
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">Cooldown (seconds)</label>
                  <input
                    type="number"
                    name="cooldown_sec"
                    value={questForm.cooldown_sec}
                    onChange={handleQuestFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                  />
                </div>
              </div>
              
              <div className="flex items-center">
                <input
                  type="checkbox"
                  name="streak_enabled"
                  checked={questForm.streak_enabled}
                  onChange={handleQuestFormChange}
                  className="h-4 w-4 rounded border-slate-700 bg-slate-800 text-violet-600 focus:ring-violet-500"
                />
                <label className="ml-2 block text-sm">Enable Streak</label>
              </div>
              
              <div className="pt-4">
                <button type="submit" className="btn-primary">
                  Create Quest
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Create Shop Item Form */}
        {activeTab === 'shop' && (
          <div className="card">
            <h2 className="text-xl font-bold mb-4">Create New Shop Item</h2>
            <form onSubmit={handleShopSubmit} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">Name</label>
                  <input
                    type="text"
                    name="name"
                    value={shopForm.name}
                    onChange={handleShopFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                    required
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">Code</label>
                  <input
                    type="text"
                    name="code"
                    value={shopForm.code}
                    onChange={handleShopFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                    required
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">Price (points)</label>
                  <input
                    type="number"
                    name="price"
                    value={shopForm.price}
                    onChange={handleShopFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                    required
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">Stock (optional)</label>
                  <input
                    type="number"
                    name="stock"
                    value={shopForm.stock}
                    onChange={handleShopFormChange}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-1">Description</label>
                <textarea
                  name="description"
                  value={shopForm.description}
                  onChange={handleShopFormChange}
                  rows={3}
                  className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
                />
              </div>
              
              <div className="flex items-center">
                <input
                  type="checkbox"
                  name="is_active"
                  checked={shopForm.is_active}
                  onChange={handleShopFormChange}
                  className="h-4 w-4 rounded border-slate-700 bg-slate-800 text-violet-600 focus:ring-violet-500"
                />
                <label className="ml-2 block text-sm">Active</label>
              </div>
              
              <div className="pt-4">
                <button type="submit" className="btn-primary">
                  Create Shop Item
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Mini Metrics */}
        <div className="card mt-6">
          <h2 className="text-xl font-bold mb-4">Mini Metrics</h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="border border-slate-700 rounded-2xl p-4">
              <div className="text-2xl font-bold text-violet-400">24</div>
              <div className="text-slate-400">Total Quests</div>
            </div>
            <div className="border border-slate-700 rounded-2xl p-4">
              <div className="text-2xl font-bold text-violet-400">12</div>
              <div className="text-slate-400">Active Quests</div>
            </div>
            <div className="border border-slate-700 rounded-2xl p-4">
              <div className="text-2xl font-bold text-violet-400">8</div>
              <div className="text-slate-400">Shop Items</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Admin;