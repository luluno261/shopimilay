import { useState } from 'react'
import GrowthOverview from './GrowthOverview'
import GrowthQuickActions from './GrowthQuickActions'
import AutomationTab from './AutomationTab'
import SegmentationTab from './SegmentationTab'
import CaptureTab from './CaptureTab'
import AdsTab from './AdsTab'

type TabType = 'overview' | 'automation' | 'segmentation' | 'capture' | 'ads'

export default function GrowthCommandCenter() {
  const [activeTab, setActiveTab] = useState<TabType>('overview')

  const tabs = [
    { id: 'overview' as TabType, label: 'Vue d\'ensemble', icon: '📊' },
    { id: 'automation' as TabType, label: 'Automations', icon: '🤖' },
    { id: 'segmentation' as TabType, label: 'Segmentation', icon: '🎯' },
    { id: 'capture' as TabType, label: 'Capture', icon: '📥' },
    { id: 'ads' as TabType, label: 'Publicités', icon: '📢' },
  ]

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      {/* Tabs Navigation */}
      <div className="border-b border-gray-200">
        <nav className="flex space-x-1 px-4" aria-label="Tabs">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                px-4 py-3 text-sm font-medium border-b-2 transition-colors
                ${
                  activeTab === tab.id
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }
              `}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="p-6">
        {activeTab === 'overview' && (
          <div>
            <GrowthQuickActions />
            <div className="mt-8">
              <GrowthOverview />
            </div>
          </div>
        )}
        {activeTab === 'automation' && <AutomationTab />}
        {activeTab === 'segmentation' && <SegmentationTab />}
        {activeTab === 'capture' && <CaptureTab />}
        {activeTab === 'ads' && <AdsTab />}
      </div>
    </div>
  )
}

