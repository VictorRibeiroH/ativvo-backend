# Integração com Tabela TACO

## 📋 O que é a Tabela TACO?

A **Tabela TACO** (Tabela Brasileira de Composição de Alimentos) é mantida pela UNICAMP e contém informações nutricionais detalhadas de alimentos típicos brasileiros.

## 🚀 Como usar

### Backend (Go)

#### 1. Buscar alimentos na TACO

**Endpoint:** `GET /api/foods/taco/search?q={query}`

**Headers:**
```
Authorization: Bearer {token}
```

**Exemplo:**
```bash
curl -X GET "http://localhost:8080/api/foods/taco/search?q=arroz" \
  -H "Authorization: Bearer seu_token_aqui"
```

**Resposta:**
```json
{
  "foods": [
    {
      "id": 123,
      "name": "Arroz, integral, cozido",
      "category": "Cereais e derivados",
      "calories": 123.5,
      "protein": 2.6,
      "carbs": 25.8,
      "fat": 1.0,
      "fiber": 2.7,
      "serving_size": 100,
      "taco_id": 123,
      "source": "TACO"
    }
  ],
  "total": 1,
  "source": "TACO"
}
```

#### 2. Importar alimento da TACO

**Endpoint:** `POST /api/foods/taco/import`

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Body:**
```json
{
  "taco_id": 123
}
```

**Resposta:**
```json
{
  "message": "Food imported successfully from TACO",
  "food": {
    "id": "uuid-gerado",
    "name": "Arroz, integral, cozido",
    "calories": 123.5,
    "protein": 2.6,
    "carbs": 25.8,
    "fat": 1.0,
    "serving_size": 100,
    "created_by_id": "user-uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Frontend (React/Next.js)

#### 1. Usar o componente de busca

```tsx
import { TACOFoodSearch } from '@/components/taco-food-search'

function MyPage() {
  const handleFoodImported = (food) => {
    console.log('Alimento importado:', food)
    // Atualizar lista de alimentos, etc.
  }

  return (
    <div>
      <TACOFoodSearch onFoodImported={handleFoodImported} />
    </div>
  )
}
```

#### 2. Usar as funções da API diretamente

```typescript
import { searchTACOFoods, importTACOFood } from '@/services/api'

// Buscar alimentos
const result = await searchTACOFoods('frango')
console.log(result.foods)

// Importar um alimento
const importedFood = await importTACOFood(123)
console.log(importedFood)
```

## 🔧 Configuração

### API da TACO

O código usa uma API wrapper da TACO hospedada em:
```
https://api-taco.herokuapp.com/api
```

Se quiser usar outra fonte de dados da TACO, altere a constante `TACO_BASE_URL` em:
```
ativvo-backend/internal/services/taco.go
```

### Alternativas de API TACO:

1. **API oficial UNICAMP:** `http://tbca-v4.cfbr.org.br`
2. **API wrapper Heroku:** `https://api-taco.herokuapp.com/api`
3. **Criar sua própria API** fazendo scraping do site oficial
4. **Baixar CSV da TACO** e criar tabela local no PostgreSQL

## 📦 Estrutura de arquivos criados

```
ativvo-backend/
├── internal/
│   ├── services/
│   │   └── taco.go              # Serviço de integração TACO
│   └── handlers/
│       └── food.go              # Handlers SearchTACOFoods e ImportTACOFood
└── main.go                      # Rotas /foods/taco/search e /foods/taco/import

ativvo-frontend/
├── components/
│   └── taco-food-search.tsx     # Componente de busca TACO
└── services/
    └── api.ts                   # Funções searchTACOFoods e importTACOFood
```

## 🎯 Fluxo de uso

1. Usuário digita "arroz" no campo de busca
2. Frontend chama `searchTACOFoods("arroz")`
3. Backend faz requisição para API TACO
4. Backend retorna lista de alimentos encontrados
5. Usuário clica em "Importar" em um alimento
6. Frontend chama `importTACOFood(tacoId)`
7. Backend busca detalhes do alimento na TACO
8. Backend cria alimento no banco de dados local
9. Frontend atualiza a lista de alimentos

## 🐛 Troubleshooting

### Erro: "Failed to search TACO foods"

- Verifique se a API TACO está online
- Teste manualmente: `curl https://api-taco.herokuapp.com/api/food/search?q=arroz`
- Considere usar cache local ou fallback

### Erro: "No authentication token found"

- Certifique-se de estar logado
- Verifique se o token está no localStorage
- Faça login novamente se necessário

### Alimentos com informações incompletas

- A TACO pode não ter todos os nutrientes para todos os alimentos
- Valores podem aparecer como 0 ou vazios
- É normal para alguns alimentos específicos

## 📊 Dados da TACO

A Tabela TACO contém aproximadamente **600 alimentos** típicos brasileiros com informações sobre:

- Energia (kcal)
- Proteínas (g)
- Lipídeos/Gorduras (g)
- Carboidratos (g)
- Fibra alimentar (g)
- Vitaminas
- Minerais
- E muito mais!

## 🔄 Próximos passos

1. **Cache local:** Salvar resultados da TACO para reduzir chamadas à API
2. **Sincronização:** Atualizar alimentos TACO periodicamente
3. **Favoritos:** Permitir marcar alimentos TACO favoritos
4. **Filtros:** Adicionar filtros por categoria, macros, etc.
5. **Offline:** Baixar toda tabela TACO para uso offline
