export function useDraftCache() {
  const KEY_PREFIX = 'ocean_survey_draft_'

  function getKey(pointCloudId: string, userId: string): string {
    return `${KEY_PREFIX}${userId}_${pointCloudId}`
  }

  function getDraft(pointCloudId: string, userId: string): any | null {
    try {
      const key = getKey(pointCloudId, userId)
      const value = localStorage.getItem(key)
      return value ? JSON.parse(value) : null
    } catch (e) {
      return null
    }
  }

  function setDraft(pointCloudId: string, userId: string, data: any): void {
    try {
      const key = getKey(pointCloudId, userId)
      localStorage.setItem(key, JSON.stringify(data))
    } catch (e) {
      console.error('Failed to save draft:', e)
    }
  }

  function clearDraft(pointCloudId: string, userId: string): void {
    const key = getKey(pointCloudId, userId)
    localStorage.removeItem(key)
  }

  function getAllDrafts(userId: string): { pointCloudId: string; data: any }[] {
    const drafts: { pointCloudId: string; data: any }[] = []
    const prefix = `${KEY_PREFIX}${userId}_`
    try {
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i)
        if (key && key.startsWith(prefix)) {
          const pointCloudId = key.substring(prefix.length)
          const value = localStorage.getItem(key)
          if (value) {
            drafts.push({
              pointCloudId,
              data: JSON.parse(value)
            })
          }
        }
      }
    } catch (e) {
      console.error('Failed to get all drafts:', e)
    }
    return drafts
  }

  function clearAllDrafts(userId: string): void {
    const prefix = `${KEY_PREFIX}${userId}_`
    const keysToRemove: string[] = []
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.startsWith(prefix)) {
        keysToRemove.push(key)
      }
    }
    keysToRemove.forEach(key => localStorage.removeItem(key))
  }

  return {
    getDraft,
    setDraft,
    clearDraft,
    getAllDrafts,
    clearAllDrafts
  }
}
