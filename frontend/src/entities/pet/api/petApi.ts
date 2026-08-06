import { baseApi } from '@/shared/api/baseApi'
import type { PetDto, PetResponse } from '../model/types'

const PET_URL = 'http://localhost:8082/api/v1/pet'

export const petApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getPet: builder.query<PetResponse, void>({
      query: () => ({ url: PET_URL }),
      providesTags: ['Pet'],
    }),
    createPet: builder.mutation<PetResponse, PetDto>({
      query: (petData) => ({ url: PET_URL, method: 'POST', body: petData }),
      invalidatesTags: ['Pet'],
    }),
    feedPet: builder.mutation<PetResponse, void>({
      query: () => ({ url: `${PET_URL}/feed`, method: 'POST' }),
      invalidatesTags: ['Pet'],
    }),
    strokePet: builder.mutation<PetResponse, void>({
      query: () => ({ url: `${PET_URL}/stroke`, method: 'POST' }),
      invalidatesTags: ['Pet'],
    }),
  }),
})

export const {
  useLazyGetPetQuery,
  useCreatePetMutation,
  useFeedPetMutation,
  useStrokePetMutation,
} = petApi
