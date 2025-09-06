import React, { useState } from 'react';
import { Quest } from '../../types/quest';
import { calculateReward } from '../../utils/helpers';
import Badge from '../ui/Badge';
import CompletionDialog from './CompletionDialog';

interface QuestRowProps {
  quest: Quest;
  onComplete: (questId: string, completion: any) => void;
}

const QuestRow: React.FC<QuestRowProps> = ({ quest, onComplete }) => {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  
  const handleComplete = () => {
    if (quest.mode === 'BINARY') {
      // For binary quests, complete immediately
      onComplete(quest.id, { idempotency_key: crypto.randomUUID() });
    } else {
      // For other modes, open dialog
      setIsDialogOpen(true);
    }
  };

  const handleDialogSubmit = (completion: any) => {
    onComplete(quest.id, completion);
    setIsDialogOpen(false);
  };

  const getModeLabel = () => {
    switch (quest.mode) {
      case 'BINARY': return 'Complete';
      case 'PARTIAL': return 'Log Progress';
      case 'PER_MINUTE': return 'Log Time';
      default: return 'Complete';
    }
  };

  const getRewardHint = () => {
    return calculateReward(quest);
  };

  return (
    <>
      <div className="card">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <h3 className="font-medium">{quest.title}</h3>
              <Badge variant="difficulty" type={quest.difficulty} />
              <Badge variant="mode" type={quest.mode} />
            </div>
            <p className="text-sm text-slate-400 mb-3">{quest.description}</p>
            <div className="text-sm">
              <span className="text-slate-300">Reward: </span>
              <span className="font-medium text-violet-400">{getRewardHint()}</span>
            </div>
          </div>
          
          <button
            onClick={handleComplete}
            className="btn-primary whitespace-nowrap"
          >
            {getModeLabel()}
          </button>
        </div>
      </div>

      <CompletionDialog
        quest={quest}
        isOpen={isDialogOpen}
        onClose={() => setIsDialogOpen(false)}
        onSubmit={handleDialogSubmit}
      />
    </>
  );
};

export default QuestRow;