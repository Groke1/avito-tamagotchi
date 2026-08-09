import { baseApi } from '@/shared/api/baseApi'
import type { CompleteTaskResponse, TaskResponse, TasksResponse } from '../model/types'

const TASKS_URL = import.meta.env.VITE_API_TASKS_URL || 'http://localhost:8081/api/v1/tasks'

export const taskApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getTasks: builder.query<TasksResponse, void>({
      query: () => ({ url: TASKS_URL }),
      providesTags: ['Tasks'],
    }),
    getTask: builder.query<TaskResponse, string>({
      query: (id) => ({ url: `${TASKS_URL}/${id}` }),
      providesTags: ['Tasks'],
    }),
    completeTask: builder.mutation<CompleteTaskResponse, string>({
      query: (id) => ({ url: `${TASKS_URL}/${id}/complete`, method: 'PUT' }),
      invalidatesTags: ['Tasks', 'User', 'Pet'],
    }),
  }),
})

export const { useGetTasksQuery, useGetTaskQuery, useCompleteTaskMutation } = taskApi
