package com.example.farmscurity.service

import com.example.farmscurity.model.ApiResponse
import com.example.farmscurity.model.ApiResponseHistory
import com.example.farmscurity.model.ApiResponseSensor
import com.example.farmscurity.model.History
import com.example.farmscurity.model.SetActivationSensorRequest
import com.example.farmscurity.model.SetIsRead
import com.google.gson.GsonBuilder
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.create
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path

interface ApiService {
    @GET("histories/")
    suspend fun getHistory(): ApiResponseHistory

    @GET("sensors/")
    suspend fun getSensor(): ApiResponseSensor

    @PUT("sensor/")
    suspend fun setIsActive(@Body request: SetActivationSensorRequest): ApiResponse

    @POST("capture")
    suspend fun capture(): ApiResponse

    @POST("turn_on")
    suspend fun alarmActive(): ApiResponse

    @POST("turn_off")
    suspend fun alarmNonActive(): ApiResponse

    @PUT("history/{historyId}")
    suspend fun setIsRead(@Path("historyId") historyId: String): ApiResponse

    @DELETE("histories")
    suspend fun deleteHistory(): ApiResponse
}

object RetrofitInstance {
    private const val BASE_URL = "http://farm.test.dihara.my.id/api/"
//    private const val BASE_URL = "http://192.168.1.7:8080/api/"

    private val gson = GsonBuilder()
        .setDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'")
        .create()

    val api: ApiService by lazy {
        Retrofit.Builder()
            .baseUrl(BASE_URL)
            .addConverterFactory(GsonConverterFactory.create(gson))
            .build()
            .create(ApiService::class.java)
    }
}
