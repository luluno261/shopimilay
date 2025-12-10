import React from 'react'

interface ChartProps {
  title: string
  data: Array<{ label: string; value: number }>
  type?: 'line' | 'bar' | 'pie'
}

export default function Chart({ title, data, type = 'line' }: ChartProps) {
  // Composant de graphique simple - peut être remplacé par une bibliothèque comme Chart.js ou Recharts
  const maxValue = Math.max(...data.map(d => d.value), 1)
  
  return (
    <div className="bg-white rounded-lg shadow-md p-6 border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">{title}</h3>
      <div className="h-64 flex items-end justify-between gap-2">
        {data.map((item, index) => (
          <div key={index} className="flex-1 flex flex-col items-center">
            <div
              className="w-full bg-blue-500 rounded-t transition-all duration-300 hover:bg-blue-600"
              style={{ height: `${(item.value / maxValue) * 100}%` }}
              title={`${item.label}: ${item.value}`}
            />
            <span className="text-xs text-gray-600 mt-2 truncate w-full text-center">
              {item.label}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}

