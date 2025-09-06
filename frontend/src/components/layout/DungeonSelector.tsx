import React, { useState } from 'react';
import { Dungeon } from '../../types/dungeon';

interface DungeonSelectorProps {
  dungeons: Dungeon[];
  selectedDungeon: Dungeon | null;
  onDungeonChange: (dungeonId: string) => void;
}

const DungeonSelector: React.FC<DungeonSelectorProps> = ({ 
  dungeons, 
  selectedDungeon, 
  onDungeonChange 
}) => {
  const [isOpen, setIsOpen] = useState(false);

  const handleSelect = (dungeonId: string) => {
    onDungeonChange(dungeonId);
    setIsOpen(false);
  };

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="btn-secondary flex items-center space-x-2"
      >
        <span>{selectedDungeon?.title || 'Select Dungeon'}</span>
        <svg 
          xmlns="http://www.w3.org/2000/svg" 
          width="16" 
          height="16" 
          viewBox="0 0 24 24" 
          fill="none" 
          stroke="currentColor" 
          strokeWidth="2" 
          strokeLinecap="round" 
          strokeLinejoin="round"
          className={`transition-transform ${isOpen ? 'rotate-180' : ''}`}
        >
          <path d="m6 9 6 6 6-6" />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-64 rounded-2xl border border-slate-700 bg-slate-800 shadow-lg z-20">
          <div className="py-1">
            {dungeons.map((dungeon) => (
              <button
                key={dungeon.id}
                onClick={() => handleSelect(dungeon.id)}
                className={`w-full px-4 py-2 text-left text-sm hover:bg-slate-700 ${
                  selectedDungeon?.id === dungeon.id ? 'bg-slate-700' : ''
                }`}
              >
                {dungeon.title}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default DungeonSelector;