package com.example.farmscurity.screen

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavController
import com.example.farmscurity.model.ApiServiceViewModel
import com.example.farmscurity.service.NetworkImage
import com.example.farmscurity.ui.theme.MyTypography

@Composable
fun Capture(
    navController: NavController,
    pictureId: String?,
    isCapture: Boolean?,
    viewModel: ApiServiceViewModel = viewModel()
) {
    LaunchedEffect(Unit) {
        viewModel.fetchHistory()
    }

    val histories by viewModel.historyList.collectAsState()
    val matchedHistory = pictureId?.let { id ->
        histories.firstOrNull { it.pictureId == id }
    }

    Scaffold(
        modifier = Modifier.fillMaxSize(),
        topBar = {
            MyTopBar(
                title = if (isCapture == true) "Hasil Tangkap Gambar" else "Gerakan Terdeteksi",
                navController = navController
            )
        }
    ) { innerPadding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(innerPadding)
                .padding(horizontal = 32.dp, vertical = 8.dp),
        ) {
            if (pictureId != null) {
                NetworkImage(fileName = pictureId)

                if (isCapture == false) {
                    matchedHistory?.let { history ->
                        Spacer(modifier = Modifier.height(24.dp))
                        Text(
                            text = history.operation,
                            style = MyTypography.titleLarge,
                            color = Color(0xFF0F172B)
                        )
                        Text(
                            text = history.createdAt,
                            style = MyTypography.bodySmall,
                            color = Color.Gray
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            text = history.description,
                            style = MyTypography.titleMedium
                        )

                    } ?: Text("Data tidak ditemukan.", color = Color.Red)
                }

            } else {
                CircularProgressIndicator()
                Text("Memuat gambar...", modifier = Modifier.padding(top = 8.dp))
            }
        }
    }
}


@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MyTopBar(title: String, navController: NavController) {
    TopAppBar(
        title = {
            Text(
                text = title,
                color = Color(0xFF0F172B),
                fontSize = 18.sp
            )
        },
        navigationIcon = {
            IconButton(onClick = { navController.popBackStack() }) {
                Icon(
                    imageVector = Icons.Default.ArrowBack,
                    contentDescription = "Kembali",
                    tint = Color(0xFF0F172B)
                )
            }
        },
        colors = TopAppBarDefaults.topAppBarColors(
            containerColor = Color(0xFFF5F7FF)
        )
    )
}
