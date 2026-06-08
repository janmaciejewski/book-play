export function useFacilityTypes() {
  const typeLabels: Record<string, string> = {
    FOOTBALL: 'Piłka nożna',
    BASKETBALL: 'Koszykówka',
    TENNIS: 'Tenis',
    VOLLEYBALL: 'Siatkówka',
    SWIMMING: 'Pływanie',
    OTHER: 'Inne',
  }

  const typeBadgeClasses: Record<string, string> = {
    FOOTBALL:
      'text-green-700 dark:text-green-400 bg-green-100 dark:bg-green-900/30 border-green-200 dark:border-green-700/50',
    BASKETBALL:
      'text-orange-700 dark:text-orange-400 bg-orange-100 dark:bg-orange-900/30 border-orange-200 dark:border-orange-700/50',
    TENNIS:
      'text-yellow-700 dark:text-yellow-400 bg-yellow-100 dark:bg-yellow-900/30 border-yellow-200 dark:border-yellow-700/50',
    VOLLEYBALL:
      'text-blue-700 dark:text-blue-400 bg-blue-100 dark:bg-blue-900/30 border-blue-200 dark:border-blue-700/50',
    SWIMMING:
      'text-cyan-700 dark:text-cyan-400 bg-cyan-100 dark:bg-cyan-900/30 border-cyan-200 dark:border-cyan-700/50',
    OTHER:
      'text-gray-700 dark:text-gray-400 bg-gray-100 dark:bg-gray-700/30 border-gray-200 dark:border-gray-600/50',
  }

  const typeImages: Record<string, string> = {
    FOOTBALL: 'https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=800&h=600&fit=crop',
    BASKETBALL: 'https://images.unsplash.com/photo-1546519638-68e109498ffc?w=800&h=600&fit=crop',
    TENNIS: 'https://images.unsplash.com/photo-1554068865-24cecd4e34b8?w=800&h=600&fit=crop',
    VOLLEYBALL: 'https://images.unsplash.com/photo-1612872087720-bb876e2e67d1?w=800&h=600&fit=crop',
    SWIMMING: 'https://images.unsplash.com/photo-1530549387789-4c1017266635?w=800&h=600&fit=crop',
    OTHER: 'https://images.unsplash.com/photo-1770064319607-f869d94d4ec3?w=800&h=600&fit=crop',
  }

  function getTypeName(type: string): string {
    return typeLabels[type] || type
  }

  function getTypeBadgeClass(type: string): string {
    return typeBadgeClasses[type] || typeBadgeClasses.OTHER
  }

  function getTypeImage(type: string): string {
    return typeImages[type] || typeImages.OTHER
  }

  const typeOptions = [
    { value: 'FOOTBALL', label: 'Piłka nożna' },
    { value: 'BASKETBALL', label: 'Koszykówka' },
    { value: 'TENNIS', label: 'Tenis' },
    { value: 'VOLLEYBALL', label: 'Siatkówka' },
    { value: 'SWIMMING', label: 'Pływanie' },
    { value: 'OTHER', label: 'Inne' },
  ]

  return {
    typeLabels,
    typeBadgeClasses,
    typeImages,
    typeOptions,
    getTypeName,
    getTypeBadgeClass,
    getTypeImage,
  }
}