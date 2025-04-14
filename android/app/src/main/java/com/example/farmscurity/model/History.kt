package com.example.farmscurity.model

import com.google.gson.annotations.SerializedName
import java.util.Date

data class ApiResponseHistory(
    val code: Int,
    val status: String,
    val data: List<History>
)

data class History(
    @SerializedName("id")
    val id: String,
    @SerializedName("operation")
    val operation: String,
    @SerializedName("picture_id")
    val pictureId: String,
    @SerializedName("created_at")
    val createdAt: String,
    @SerializedName("description")
    val description: String,
)
