import { baseApi } from '@/shared/api/baseApi'
import type { PetDto, PetResponse, PetTicketResponse, PetTripResponse } from '../model/types'

const PET_URL = import.meta.env.VITE_API_PET_URL || 'http://localhost:8082/api/v1/pet'

export const petApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getPet: builder.query<PetResponse, void>({
      query: () => ({ url: PET_URL }),
      providesTags: ['Pet'],
    }),
    getPetTripLast: builder.query<PetTripResponse, void>({
      query: () => ({ url: `${PET_URL}/trip/last` }),
      providesTags: ['Pet'],
    }),
    createPet: builder.mutation<PetResponse, PetDto>({
      query: (petData) => ({ url: PET_URL, method: 'POST', body: petData }),
      invalidatesTags: ['Pet'],
    }),
    feedPet: builder.mutation<PetResponse, void>({
      query: () => ({ url: `${PET_URL}/feed`, method: 'POST' }),
      invalidatesTags: ['Pet', 'User'],
    }),
    strokePet: builder.mutation<PetResponse, void>({
      query: () => ({ url: `${PET_URL}/stroke`, method: 'POST' }),
      invalidatesTags: ['Pet'],
    }),
    makeTrip: builder.mutation<void, number>({
      query: (petId) => ({ url: `${PET_URL}/trip/${petId}`, method: 'POST' }),
      invalidatesTags: ['Pet', 'User', 'Leaderboard'],
    }),
    getWsTicket: builder.mutation<PetTicketResponse, void>({
      query: () => ({ url: `${PET_URL}/ws-ticket`, method: 'POST' }),
    }),
    tripPet: builder.mutation<void, number>({
      query: (petId) => ({
        url: `${PET_URL}/trip/${petId}`,
        method: 'POST',
      }),
      invalidatesTags: ['Pet'],
    }),
  }),
})

export const {
  useGetPetQuery,
  useLazyGetPetQuery,
  useGetPetTripLastQuery,
  useCreatePetMutation,
  useFeedPetMutation,
  useStrokePetMutation,
  useGetWsTicketMutation,
  useTripPetMutation,
} = petApi
