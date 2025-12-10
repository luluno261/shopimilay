import React, { useState } from 'react'

interface DiscountCodeProps {
  value: string
  onChange: (value: string) => void
  onApply: () => void
  applied: boolean
}

export default function DiscountCode({
  value,
  onChange,
  onApply,
  applied,
}: DiscountCodeProps) {
  return (
    <div className="bg-gray-50 rounded-lg p-4">
      <h3 className="text-lg font-semibold mb-2">Code de réduction</h3>
      <div className="flex gap-2">
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="Entrez votre code"
          className="flex-1 border border-gray-300 rounded-md px-3 py-2"
          disabled={applied}
        />
        {!applied ? (
          <button
            onClick={onApply}
            className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
          >
            Appliquer
          </button>
        ) : (
          <span className="text-green-600 px-4 py-2">✓ Appliqué</span>
        )}
      </div>
    </div>
  )
}

