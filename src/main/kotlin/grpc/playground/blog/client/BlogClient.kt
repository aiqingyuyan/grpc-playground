package grpc.playground.blog.client

import com.example.protobuf.blog.BlogServiceGrpcKt
import com.example.protobuf.blog.GetBlogRequest
import com.example.protobuf.blog.GetBlogResponse
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.asCoroutineDispatcher
import kotlinx.coroutines.asExecutor
import kotlinx.coroutines.runBlocking
import java.io.Closeable
import java.util.concurrent.Executors

class BlogClient private constructor(
	private val channel: ManagedChannel
) : Closeable {

	private val stub = BlogServiceGrpcKt.BlogServiceCoroutineStub(channel)

	constructor(
		channelBuilder: ManagedChannelBuilder<*>,
		dispatcher: CoroutineDispatcher
	) : this(
		channelBuilder
			.executor(dispatcher.asExecutor())
			.build()
	)

	fun getBlog(request: GetBlogRequest): GetBlogResponse = runBlocking {
		stub.getBlog(request)
	}

	override fun close() {
		channel.shutdown()
	}
}

fun main() {
	Executors
		.newFixedThreadPool(2)
		.asCoroutineDispatcher()
		.use { dispatcher ->
			BlogClient(
				ManagedChannelBuilder
					.forAddress("localhost", 50051)
					.usePlaintext(),
				dispatcher
			).use { client ->
				println("Send request 1 ...")
				var response = client.getBlog(
					GetBlogRequest.newBuilder()
						.setId(1)
						.build()
				)
				println("Received response 1")
				println(response)

				println("Send request 2 ...")
				response = client.getBlog(
					GetBlogRequest.newBuilder()
						.setId(2)
						.build()
				)
				println("Received response 2")
				println(response)
			}
		}
}
