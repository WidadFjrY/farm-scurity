package com.example.farmscurity.model

import com.google.gson.annotations.SerializedName

data class ApiResponseSensor(
    val code: Int,
    val status: String,
    val data: List<Sensor>
)

data class Sensor(
    @SerializedName("id")
    val id: String,
    @SerializedName("location")
    val location: String,
    @SerializedName("is_active")
    val isActive: Boolean,
)

data class SetActivationSensorRequest(
    @SerializedName("id")
    val id: String,
    @SerializedName("is_active")
    val isActive: Boolean,
)