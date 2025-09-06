import React, { useState } from 'react';
import { Quest } from '../../types/quest';
import { generateId } from '../../utils/helpers';

interface CompletionDialogProps {
  quest: Quest;
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (completion: any) => void;
}

const CompletionDialog: React.FC<CompletionDialogProps> = ({ 
  quest, 
  isOpen, 
  onClose, 
  onSubmit 
}) => {
  const [sliderValue, setSliderValue] = useState(50);
  const [minutesValue, setMinutesValue] = useState(30);

  if (!isOpen) return null;

  const handleSubmit = () => {
    const idempotencyKey = generateId();
    
    switch (quest.mode) {
      case 'PARTIAL':
        onSubmit({
          completion_ratio: sliderValue / 100,
          idempotency_key: idempotencyKey
        });
        break;
      case 'PER_MINUTE':
        onSubmit({
          minutes: minutesValue,
          idempotency_key: idempotencyKey
        });
        break;
      default:
        onSubmit({
          idempotency_key: idempotencyKey
        });
    }
  };

  const renderDialogContent = () => {
    switch (quest.mode) {
      case 'PARTIAL':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">
                Progress: {sliderValue}%
              </label>
              <input
                type="range"
                min="0"
                max="100"
                step="5"
                value={sliderValue}
                onChange={(e) => setSliderValue(parseInt(e.target.value))}
                className="w-full h-2 bg-slate-700 rounded-lg appearance-none cursor-pointer"
              />
              <div className="flex justify-between text-xs text-slate-400 mt-1">
                <span>0%</span>
                <span>100%</span>
              </div>
            </div>
            <div className="text-sm text-slate-300">
              You'll earn approximately {Math.round((parseFloat(quest.points_award) || 0) * (sliderValue / 100))} points
            </div>
          </div>
        );
      
      case 'PER_MINUTE':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">
                Minutes spent
              </label>
              <input
                type="number"
                min="0"
                value={minutesValue}
                onChange={(e) => setMinutesValue(parseInt(e.target.value) || 0)}
                className="w-full rounded-lg border border-slate-700 bg-slate-800 px-3 py-2"
              />
              {quest.min_minutes && (
                <p className="text-xs text-slate-400 mt-1">
                  Minimum: {quest.min_minutes} minutes
                </p>
              )}
              {quest.max_minutes && (
                <p className="text-xs text-slate-400 mt-1">
                  Maximum: {quest.max_minutes} minutes
                </p>
              )}
            </div>
            <div className="text-sm text-slate-300">
              {quest.rate_points_per_min && (
                <p>
                  You'll earn approximately {Math.round((parseFloat(quest.rate_points_per_min) || 0) * minutesValue)} points
                </p>
              )}
              {quest.max_minutes && quest.rate_points_per_min && (
                <p className="text-xs text-slate-400 mt-1">
                  Max points: {Math.round((parseFloat(quest.rate_points_per_min) || 0) * quest.max_minutes)} points
                </p>
              )}
            </div>
          </div>
        );
      
      default:
        return (
          <p>Are you sure you want to complete this quest?</p>
        );
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="card w-full max-w-md">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-bold">{quest.title}</h3>
          <button 
            onClick={onClose}
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
        
        {renderDialogContent()}
        
        <div className="flex justify-end space-x-3 mt-6">
          <button
            onClick={onClose}
            className="btn-secondary"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            className="btn-primary"
          >
            Submit
          </button>
        </div>
      </div>
    </div>
  );
};

export default CompletionDialog;