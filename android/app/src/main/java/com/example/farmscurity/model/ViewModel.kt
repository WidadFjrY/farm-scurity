package com.example.farmscurity.model

import android.util.Log
import androidx.compose.runtime.Composable
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.farmscurity.service.RetrofitInstance
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class ApiServiceViewModel : ViewModel() {
    private val _historyList = MutableStateFlow<List<History>>(emptyList())
    val historyList: StateFlow<List<History>> = _historyList

    private val _sensorList = MutableStateFlow<List<Sensor>>(emptyList())
    val sensorList: StateFlow<List<Sensor>> = _sensorList

    private val _isActive = MutableStateFlow<List<ApiResponse>>(emptyList())
    val isActive: StateFlow<List<ApiResponse>> = _isActive

    private val _capture = MutableStateFlow<List<ApiResponse>>(emptyList())
    val capture: StateFlow<List<ApiResponse>> = _capture

    private val _alarm = MutableStateFlow<List<ApiResponse>>(emptyList())
    val alarm: StateFlow<List<ApiResponse>> = _alarm

    init {
        fetchHistory()
        fetchSensor()
    }

    fun fetchSensor() {
        viewModelScope.launch {
            try {
                val response = RetrofitInstance.api.getSensor()

                if (response.code == 200) {
                    _sensorList.value = response.data
                } else {
                    _sensorList.value = emptyList()
                }
            } catch (e: Exception) {
                _sensorList.value = emptyList()
            }
        }
    }

    fun fetchHistory() {
        viewModelScope.launch {
            try {
                val response = RetrofitInstance.api.getHistory()

                if (response.code == 200) {
                    _historyList.value = response.data
                    if (response.data.isEmpty()) {
                        _historyList.value = emptyList()
                    }
                } else {
                    _historyList.value = emptyList()
                }
            } catch (e: Exception) {
                _historyList.value = emptyList()
            }
        }
    }

    fun setIsActive(sensorId: String, isActive: Boolean, onResult: (Boolean) -> Unit) {
        viewModelScope.launch {
            try {
                val request = SetActivationSensorRequest(sensorId, isActive)
                val response = RetrofitInstance.api.setIsActive(request)
                if (response.code == 200) {
                    onResult(true)
                    _isActive.value = listOf(response)
                    fetchSensor()
                } else {
                    onResult(false)
                    _isActive.value = emptyList()
                }
            } catch (e: Exception) {
                _isActive.value = emptyList()
                onResult(false)
            }
        }
    }

    suspend fun capture(): Boolean {
        return try {
            val response = RetrofitInstance.api.capture()
            if (response.code == 200) {
                _capture.value = listOf(response)
                true
            } else {
                _capture.value = emptyList()
                false
            }
        } catch (e: Exception) {
            _capture.value = emptyList()
            false
        }
    }

    fun setIsRead(historyId: String){
        viewModelScope.launch {
            try {
                val response = RetrofitInstance.api.setIsRead(historyId)
                if (response.code == 200){
                    //
                }
            } catch (e: Exception){
                //
            }
        }
    }

    fun alarm(isActive: Boolean, onResult: (Boolean) -> Unit) {
        viewModelScope.launch {
            try {
                if (isActive) {
                    val response = RetrofitInstance.api.alarmActive()
                    if (response.code == 200) {
                        onResult(true)
                        _alarm.value = listOf(response)
                    } else {
                        onResult(false)
                        _alarm.value = emptyList()
                    }
                } else {
                    val response = RetrofitInstance.api.alarmNonActive()
                    if (response.code == 200) {
                        onResult(true)
                        _alarm.value = listOf(response)
                    } else {
                        onResult(false)
                        _alarm.value = emptyList()
                    }
                }
            } catch (e: Exception) {
                onResult(false)
                _alarm.value = emptyList()
            }
        }
    }
}