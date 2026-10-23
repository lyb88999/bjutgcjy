import service from '@/utils/request'

export const getExpertAgentStatus = ({ signal } = {}) => {
  return service({
    url: '/expertDatabase/agent/status',
    method: 'get',
    timeout: 15000,
    signal,
    donNotShowLoading: true,
    donNotShowError: true
  })
}

export const chatWithExpertAgent = (data, { signal } = {}) => {
  return service({
    url: '/expertDatabase/agent/chat',
    method: 'post',
    data,
    timeout: 90000,
    signal,
    donNotShowLoading: true,
    donNotShowError: true
  })
}
